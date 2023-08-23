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
	WorkerInitStatus WorkerStatus = iota
	WorkerRunningStatus
	WorkerIdleStatus
	WorkerOfflineStatus
	WorkerFailedStatus
)

type Worker struct {
	Name              string
	Address           string
	CreatedAt         time.Time
	LastHeartbeatTime time.Time
	Status            WorkerStatus
}

func NewWorker(name string, address string) *Worker {
	return &Worker{
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
		Status:    WorkerInitStatus,
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

func (s *Worker) SetRunningStatus() {
	s.Status = WorkerRunningStatus
}

func (s *Worker) SetIdleStatus() {
	s.Status = WorkerIdleStatus
}

func (s *Worker) SetOfflineStatus() {
	s.Status = WorkerOfflineStatus
}

func (s *Worker) SetFailedStatus() {
	s.Status = WorkerFailedStatus
}
