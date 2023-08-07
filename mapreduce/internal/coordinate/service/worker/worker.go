package worker

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/log"
	"sync"
	"time"
)

const (
	MapHeartbeatTimeout = 1 * time.Minute
)

var (
	workers sync.Map
)

type Worker struct {
	Name              string
	Address           string
	CreatedAt         time.Time
	LastHeartbeatTime time.Time
	Status            WorkerStatus
}

func NewWorker(name string, address string) (*Worker, error) {
	if _, ok := workers.Load(name); ok {
		return nil, errors.New(fmt.Sprintf("map worker[%s] exists", name))
	}

	worker := Worker{
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
		Status:    InitStatus,
	}

	workers.Store(name, &worker)
	return &worker, nil
}

func GetWorker(name string) (*Worker, error) {
	worker, ok := workers.Load(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf("map worker[%s] does not exists", name))
	}

	return worker.(*Worker), nil
}

func (worker *Worker) UpdateHeartbeat() {
	log.Info(fmt.Sprintf("map worker[%s] heartbeat", worker.Name))
	worker.LastHeartbeatTime = time.Now()
}

func (worker *Worker) GetCurrentStatus() WorkerStatus {
	if worker.Status == RunningStatus {
		if time.Now().Sub(worker.LastHeartbeatTime) > MapHeartbeatTimeout {
			worker.SetOfflineStatus()
		}
	}
	return worker.Status
}

func (worker *Worker) SetRunningStatus() {
	worker.Status = RunningStatus
}

func (worker *Worker) SetIdleStatus() {
	worker.Status = IdleStatus
}

func (worker *Worker) SetOfflineStatus() {
	worker.Status = OfflineStatus
}

func (worker *Worker) SetFailedStatus() {
	worker.Status = FailedStatus
}
