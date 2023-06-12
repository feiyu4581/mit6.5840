package map_worker

import (
	"fmt"
	"github.com/pkg/errors"
	"mapreduce/internal/pkg/log"
	"sync"
	"time"
)

var (
	workers sync.Map
)

type MapWorker struct {
	Name              string
	Address           string
	CreatedAt         time.Time
	LastHeartbeatTime time.Time
	Status            WorkerStatus
}

func NewWorker(name string, address string) (*MapWorker, error) {
	if _, ok := workers.Load(name); ok {
		return nil, errors.New(fmt.Sprintf("map worker[%s] exists", name))
	}

	worker := MapWorker{
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
		Status:    InitStatus,
	}

	workers.Store(name, &worker)
	return &worker, nil
}

func GetWorker(name string) (*MapWorker, error) {
	worker, ok := workers.Load(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf("map worker[%s] does not exists", name))
	}

	return worker.(*MapWorker), nil
}

func (worker *MapWorker) UpdateHeartbeat() {
	log.Info(fmt.Sprintf("map worker[%s] heartbeat", worker.Name))
	worker.LastHeartbeatTime = time.Now()
}

func (worker *MapWorker) SetRunningStatus() {
	worker.Status = RunningStatus
}

func (worker *MapWorker) SetIdleStatus() {
	worker.Status = IdleStatus
}

func (worker *MapWorker) SetOfflineStatus() {
	worker.Status = OfflineStatus
}

func (worker *MapWorker) SetFailedStatus() {
	worker.Status = FailedStatus
}
