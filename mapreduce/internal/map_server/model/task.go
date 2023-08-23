package model

import (
	"mapreduce/internal/pkg/model"
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
	Name      string
	Filename  string
	Params    []string
	StartTime time.Time
	Function  *model.Function
	Status    Status
}

func NewTask(name string, filename string, params []string) *Task {
	return &Task{
		Name:     name,
		Filename: filename,
		Params:   params,
		Status:   InitStatus,
	}
}

func (task *Task) Start() {
}

func (task *Task) SetFunction(function *model.Function) {
	task.Function = function
}

func (task *Task) GetStatus() Status {
	return task.Status
}

func (task *Task) IsFinish() bool {
	return task.Status == DoneStatus || task.Status == FailedStatus
}
