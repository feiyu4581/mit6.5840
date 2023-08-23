package coordinator

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/coordinate/model"
	"sync"
)

var (
	coordinatorService *Service
	coordinatorOnce    sync.Once
)

type Service struct {
	sync.RWMutex

	CurrentTasks map[string]*model.Task
	WorkerMaps   map[string]*model.Worker
	IdleWorker   *model.IdleWorkerManager
}

func GetService() *Service {
	coordinatorOnce.Do(func() {
		coordinatorService = &Service{
			CurrentTasks: make(map[string]*model.Task),
			IdleWorker:   model.NewIdleWorkerManager(),
		}
	})

	return coordinatorService
}

func (s *Service) NewWorker(name, address string) error {
	if _, ok := s.WorkerMaps[name]; ok {
		return errors.New(fmt.Sprintf("map worker[%s] exists", name))
	}

	wk := model.NewWorker(name, address)
	s.WorkerMaps[name] = wk
	s.IdleWorker.Add(wk)
	return nil
}

func (s *Service) GetWorker(name string) (*model.Worker, error) {
	worker, ok := s.WorkerMaps[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("map worker[%s] does not exists", name))
	}

	return worker, nil
}

func (s *Service) StoreTask(task *model.Task) {
	s.Lock()
	defer s.Unlock()

	s.CurrentTasks[task.Name] = task
}

func (s *Service) StartTask(task *model.Task) error {
	task.Status = model.RunningMapStatus
	workers, err := s.IdleWorker.ChooseWorkers(task.MNums)
	if err != nil {
		return err
	}

	task.Workers = workers

	return nil
}

func (s *Service) AddTask(task *model.Task) error {
	s.StoreTask(task)
	return s.StartTask(task)
}
