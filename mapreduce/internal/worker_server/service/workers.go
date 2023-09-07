package service

import (
	"fmt"
	"github.com/pkg/errors"
	functionModel "mapreduce/internal/pkg/model"
	"mapreduce/internal/worker_server/model"
)

var workerManager *WorkerManager

func GetWorkerManager() *WorkerManager {
	return workerManager
}

func init() {
	workerManager = &WorkerManager{
		Functions: make(map[string]*functionModel.Function),
		SingleC:   make(chan struct{}, 1),
	}
}

type TaskHistory struct {
	Task *model.Task
}

type WorkerManager struct {
	Functions     map[string]*functionModel.Function
	CurrentTask   *model.Task
	TaskHistories []TaskHistory
	SingleC       chan struct{}
	Mode          functionModel.WorkerMode
}

func (manager *WorkerManager) SetWorkMode(mode functionModel.WorkerMode) {
	manager.Mode = mode
}

func (manager *WorkerManager) Register(function *functionModel.Function) {
	manager.Functions[function.Name] = function
}

func (manager *WorkerManager) AddTask(task *model.Task, functionName string) error {
	if manager.CurrentTask != nil && !manager.CurrentTask.IsFinish() {
		return errors.New("current task is running")
	}

	if function, ok := manager.Functions[functionName]; ok {
		task.SetFunction(function)
	} else {
		return errors.New(fmt.Sprintf("function[%s] is not registered", functionName))
	}

	manager.CurrentTask = task
	manager.TaskHistories = append(manager.TaskHistories, TaskHistory{
		Task: task,
	})

	go func() {
		<-task.DoneC
		manager.SingleC <- struct{}{}
	}()

	manager.SingleC <- struct{}{}
	return task.Start(manager.Mode)
}
