package model

import (
	"time"
)

type TaskStatus int

const (
	InitStatus TaskStatus = iota
	WaitingMapStatus
	RunningMapStatus
	WaitingReduceStatus
	RunningReduceStatus
	FinishedStatus
	FailedStatus
)

type TaskRequest struct {
	FileNames []string `json:"file_names" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Worker    string   `json:"worker" binding:"required"`
	MNums     int      `json:"m_nums" binding:"required"`
	RNums     int      `json:"r_nums" binding:"required"`
}

type TaskWorker struct {
	Worker    *Worker
	StartTime time.Time
	EndTime   time.Time
	FileNames []string
	IsDone    bool
	IsSuccess bool
}

func NewTaskWorker(worker *Worker) *TaskWorker {
	return &TaskWorker{
		Worker:    worker,
		StartTime: time.Now(),
	}
}

type Task struct {
	FileNames []string   `json:"file_names"`
	TaskId    int64      `json:"task_id"`
	Name      string     `json:"name"`
	Worker    string     `json:"worker"`
	MNums     int        `json:"m_nums"`
	RNums     int        `json:"r_nums"`
	Status    TaskStatus `json:"status"`

	SingleC chan struct{} `json:"-"`

	// TODO 等待心跳上报数据，并且更改 Status 状态
	MapWorkers    map[string]*TaskWorker `json:"map_workers"`
	ReduceWorkers map[string]*TaskWorker `json:"reduce_workers"`
}

func NewTask(fileNames []string, name string, worker string, mNums int, rNums int) *Task {
	return &Task{
		Name:          name,
		FileNames:     fileNames,
		MNums:         mNums,
		RNums:         rNums,
		Worker:        worker,
		Status:        InitStatus,
		SingleC:       make(chan struct{}),
		MapWorkers:    make(map[string]*TaskWorker),
		ReduceWorkers: make(map[string]*TaskWorker),
	}
}

func (task *Task) WaitingStatus(status TaskStatus) bool {
	ticker := time.NewTicker(time.Second * 5)
	for task.Status != status {
		select {
		case <-task.SingleC:
		case <-ticker.C:
		}

		if task.Status == FailedStatus {
			return false
		}
	}

	return true
}
