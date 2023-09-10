package coordinator

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/pkg/rpc/coordinate"
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
	} else {
		s.ReduceWorkers[name] = wk
		s.IdleReduceWorker.Add(wk)
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

// TODO close 的收尾情况，释放资源，上报结果
func (s *Service) CloseTask(task *model.Task) {

}

func (s *Service) StartTask(task *model.Task) {
	defer s.CloseTask(task)

	task.Status = model.WaitingMapStatus
	// 强制等待 ，直到有足够的 worker
	workers := s.IdleMapWorker.ChooseWorkers(task.MNums)
	for _, worker := range workers {
		task.MapWorkers[worker.Name] = model.NewTaskWorker(worker)
	}

	// TODO 开始 map 任务调度
	task.Status = model.RunningMapStatus
	if !task.WaitingStatus(model.WaitingReduceStatus) {
		return
	}

	// 强制等待 worker，直到有足够的 worker
	workers = s.IdleReduceWorker.ChooseWorkers(task.RNums)
	for _, worker := range workers {
		task.ReduceWorkers[worker.Name] = model.NewTaskWorker(worker)
	}

	// TODO 开始 reduce 任务调度
	task.Status = model.RunningMapStatus
	task.WaitingStatus(model.FinishedStatus)
}

func (s *Service) AddTask(task *model.Task) {
	s.StoreTask(task)
	go s.StartTask(task)
}
