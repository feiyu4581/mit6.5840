package model

import (
	"fmt"
	"mapreduce/internal/pkg/log"
	"time"
)

type WorkerStatus int

const (
	MapHeartbeatTimeout = 1 * time.Minute
)

const (
	MapWorkerType    = "map"
	ReduceWorkerType = "reduce"
)

const (
	WorkerInitStatus WorkerStatus = iota
	WorkerRunningStatus
	WorkerOfflineStatus
)

type Worker struct {
	Name              string
	Address           string
	CreatedAt         time.Time
	LastHeartbeatTime time.Time
	Status            WorkerStatus
	WorkerType        string
}

func NewWorker(name string, address string, isMap bool) *Worker {
	workerType := ReduceWorkerType
	if isMap {
		workerType = MapWorkerType
	}
	return &Worker{
		Name:       name,
		Address:    address,
		CreatedAt:  time.Now(),
		Status:     WorkerInitStatus,
		WorkerType: workerType,
	}
}

func (s *Worker) UpdateHeartbeat() {
	log.Info(fmt.Sprintf("map s[%s] heartbeat", s.Name))
	s.LastHeartbeatTime = time.Now()
}

func (s *Worker) GetCurrentStatus() WorkerStatus {
	if s.Status == WorkerRunningStatus {
		if time.Now().Sub(s.LastHeartbeatTime) > MapHeartbeatTimeout {
			s.SetOfflineStatus()
		}
	}
	return s.Status
}

func (s *Worker) IsRunningStatus() bool {
	return s.Status == WorkerRunningStatus
}

func (s *Worker) SetRunningStatus() {
	s.Status = WorkerRunningStatus
}

func (s *Worker) SetOfflineStatus() {
	s.Status = WorkerOfflineStatus
}
