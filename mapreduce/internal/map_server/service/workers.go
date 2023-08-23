package service

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/map_server/model"
	functionModel "mapreduce/internal/pkg/model"
)

var workerManager *WorkerManager

func GetWorkerManager() *WorkerManager {
	return workerManager
}

func init() {
	workerManager = &WorkerManager{
		Functions: make(map[string]*functionModel.Function),
	}
}

type TaskHistory struct {
	Task *model.Task
}

type WorkerManager struct {
	Functions     map[string]*functionModel.Function
	CurrentTask   *model.Task
	TaskHistories []TaskHistory
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

	task.Start()
	return nil
}
