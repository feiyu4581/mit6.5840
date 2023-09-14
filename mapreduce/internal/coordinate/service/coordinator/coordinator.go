package coordinator

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/rpc/coordinate"
	"strings"
	"sync"
	"time"
)

var (
	coordinatorService *Service
	coordinatorOnce    sync.Once
)

type Service struct {
	sync.RWMutex

	CurrentTasks  map[string]*model.Task
	MapWorkers    map[string]*model.Worker
	ReduceWorkers map[string]*model.Worker

	IdleMapWorker    *model.IdleWorkerManager
	IdleReduceWorker *model.IdleWorkerManager
}

func GetService() *Service {
	coordinatorOnce.Do(func() {
		coordinatorService = &Service{
			CurrentTasks:     make(map[string]*model.Task),
			MapWorkers:       make(map[string]*model.Worker),
			ReduceWorkers:    make(map[string]*model.Worker),
			IdleReduceWorker: model.NewIdleWorkerManager(),
			IdleMapWorker:    model.NewIdleWorkerManager(),
		}
	})

	return coordinatorService
}

func (s *Service) NewWorker(name, address string, mode coordinate.ClientMode) error {
	if _, ok := s.MapWorkers[name]; ok {
		return errors.New(fmt.Sprintf("map worker[%s] exists", name))
	}

	if _, ok := s.ReduceWorkers[name]; ok {
		return errors.New(fmt.Sprintf("reduce worker[%s] exists", name))
	}

	wk := model.NewWorker(name, address, mode == coordinate.ClientMode_MapMode)
	if mode == coordinate.ClientMode_MapMode {
		s.MapWorkers[name] = wk
		s.IdleMapWorker.Add(wk)
		log.Info("注册 map 服务成功")
	} else {
		s.ReduceWorkers[name] = wk
		s.IdleReduceWorker.Add(wk)
		log.Info("注册 reduce 服务成功")
	}

	return nil
}

func (s *Service) GetWorker(name string) (*model.Worker, error) {
	worker, ok := s.MapWorkers[name]
	if !ok {
		worker, ok = s.ReduceWorkers[name]
	}

	if !ok {
		return nil, errors.New(fmt.Sprintf("worker[%s] does not exists", name))
	}

	return worker, nil
}

func (s *Service) StoreTask(task *model.Task) {
	s.Lock()
	defer s.Unlock()

	task.TaskId = time.Now().UnixNano()
	s.CurrentTasks[task.Name] = task
}

func (s *Service) GetTask(taskName string) *model.Task {
	s.Lock()
	defer s.Unlock()

	return s.CurrentTasks[taskName]
}

func (s *Service) CloseTask(task *model.Task) {
	log.Info("Task: %s is done: mapFilenames:%s, filenames: %s, err:%s", task.Name, task.CollectMapResult(), strings.Join(task.CollectReduceResult(), ","), task.Err)
	for _, worker := range task.MapWorkers {
		s.IdleMapWorker.Add(worker.Worker)
	}

	for _, worker := range task.ReduceWorkers {
		s.IdleReduceWorker.Add(worker.Worker)
	}
}

func (s *Service) StartTask(ctx context.Context, task *model.Task) {
	defer s.CloseTask(task)

	task.Status = model.WaitingMapStatus
	// 强制等待 ，直到有足够的 worker

	log.Info("开始等待足够的 map worker")
	workers := s.IdleMapWorker.ChooseWorkers(task.MNums)
	for _, worker := range workers {
		task.MapWorkers[worker.Name] = model.NewTaskWorker(worker)
	}

	task.Status = model.RunningMapStatus
	log.Info("开始分发 map 任务")
	if err := task.StartMapWorker(ctx); err != nil {
		task.Status = model.FailedStatus
		task.Err = err
		return
	}

	log.Info("分发 map 任务完成")
	if !task.WaitingStatus(model.WaitingReduceStatus) {
		return
	}

	// 强制等待 worker，直到有足够的 worker
	log.Info("开始等待足够的 reduce worker")
	workers = s.IdleReduceWorker.ChooseWorkers(task.RNums)
	for _, worker := range workers {
		task.ReduceWorkers[worker.Name] = model.NewTaskWorker(worker)
	}

	task.Status = model.RunningReduceStatus
	log.Info("开始分发 reduce 任务")
	if err := task.StartReduceWorker(ctx, task.CollectMapResult()); err != nil {
		task.Status = model.FailedStatus
		task.Err = err
		return
	}
	task.WaitingStatus(model.FinishedStatus)
}

func (s *Service) AddTask(ctx context.Context, task *model.Task) {
	s.StoreTask(task)
	go s.StartTask(ctx, task)
}
