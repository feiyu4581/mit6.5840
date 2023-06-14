package service

import (
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/worker"
)

var workerManager *WorkerManager

func GetWorkerManager() *WorkerManager {
	return workerManager
}

func init() {
	workerManager = &WorkerManager{
		Workers: make(map[string]worker.Worker),
	}
}

type TaskHistory struct {
	Task *Task
}

type WorkerManager struct {
	Workers       map[string]worker.Worker
	CurrentTask   *Task
	TaskHistories []TaskHistory
}

func (manager *WorkerManager) Register(workerName string, worker worker.Worker) {
	manager.Workers[workerName] = worker
}

func (manager *WorkerManager) AddTask(task *Task) error {
	if manager.CurrentTask != nil && !manager.CurrentTask.IsFinish() {
		return errors.New("current task is running")
	}

	manager.CurrentTask = task
	manager.TaskHistories = append(manager.TaskHistories, TaskHistory{
		Task: task,
	})

	task.Start()
	return nil
}
