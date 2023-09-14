package model

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/rpc/map_server"
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
	Err       error      `json:"err"`

	SingleC chan struct{} `json:"-"`

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

func (task *Task) CollectMapResult() [][]string {
	filenames := make([][]string, len(task.ReduceWorkers))
	for _, worker := range task.MapWorkers {
		for index, filename := range worker.FileNames {
			filenames[index] = append(filenames[index], filename)
		}
	}

	return filenames
}

func (task *Task) CollectReduceResult() []string {
	filenames := make([]string, 0, len(task.ReduceWorkers))
	for _, worker := range task.ReduceWorkers {
		filenames = append(filenames, worker.FileNames...)
	}

	return filenames
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

func (task *Task) StartMapWorker(ctx context.Context) error {
	clients := make([]*map_server.ServerClient, 0, len(task.MapWorkers))
	for _, worker := range task.MapWorkers {
		client, err := map_server.GetServerClient(worker.Worker.Address)
		if err != nil {
			return errors.New(fmt.Sprintf("Init worker server error"))
		}
		clients = append(clients, client)
	}

	for index, filename := range task.FileNames {
		err := clients[index].NewTask(ctx, task.Name, []string{filename}, task.Worker, int64(task.MNums), task.TaskId, int64(index))
		if err != nil {
			return errors.New(fmt.Sprintf("dispatch map worker error"))
		}
	}

	return nil
}

func (task *Task) StartReduceWorker(ctx context.Context, filenames [][]string) error {
	clients := make([]*map_server.ServerClient, 0, len(task.ReduceWorkers))
	for _, worker := range task.ReduceWorkers {
		client, err := map_server.GetServerClient(worker.Worker.Address)
		if err != nil {
			return errors.New(fmt.Sprintf("Init worker server error"))
		}
		clients = append(clients, client)
	}

	for index, subFilenames := range filenames {
		err := clients[index].NewTask(ctx, task.Name, subFilenames, task.Worker, int64(task.MNums), task.TaskId, int64(index))
		if err != nil {
			return errors.New(fmt.Sprintf("dispatch map worker error"))
		}
	}

	return nil
}
