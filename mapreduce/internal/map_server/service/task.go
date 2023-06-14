package service

import (
	"mapreduce/internal/pkg/worker"
	"time"
)

type Status int

const (
	InitStatus Status = iota
	RunningStatus
	DoneStatus
	FailedStatus
)

type Task struct {
	TaskId    string
	Name      string
	Filename  string
	Params    []string
	StartTime time.Time
	Status    Status
	Worker    worker.Worker
}

func NewTask(taskId string, name string, filename string, params []string) *Task {
	return &Task{
		TaskId:   taskId,
		Name:     name,
		Filename: filename,
		Params:   params,
		Status:   InitStatus,
	}
}

func (task *Task) Start() {
}

func (task *Task) GetStatus() Status {
	return task.Status
}

func (task *Task) IsFinish() bool {
	return task.Status == DoneStatus || task.Status == FailedStatus
}
