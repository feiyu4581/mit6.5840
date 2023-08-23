package model

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

type IdleWorker struct {
	Worker *Worker
	Next   *IdleWorker
}

type IdleWorkerManager struct {
	sync.Mutex

	Head *IdleWorker
	Tail *IdleWorker

	count int
}

func NewIdleWorkerManager() *IdleWorkerManager {
	return &IdleWorkerManager{}
}

func (m *IdleWorkerManager) Add(worker *Worker) {
	m.Lock()
	defer m.Unlock()

	idleWorker := &IdleWorker{
		Worker: worker,
	}

	if m.IsEmpty() {
		m.Head = idleWorker
		m.Tail = idleWorker
		return
	}

	m.Tail.Next = idleWorker
	m.Tail = idleWorker
	m.count += 1
}

func (m *IdleWorkerManager) Pop() *Worker {
	m.Lock()
	defer m.Unlock()

	if m.IsEmpty() {
		return nil
	}

	return m.PopWithoutLock()
}

func (m *IdleWorkerManager) PopWithoutLock() *Worker {
	worker := m.Head.Worker
	m.Head = m.Head.Next
	m.count -= 1
	return worker
}

func (m *IdleWorkerManager) IsEmpty() bool {
	return m.Head == nil
}

func (m *IdleWorkerManager) ChooseWorkers(nums int) ([]*Worker, error) {
	if m.count < nums {
		return nil, errors.New(fmt.Sprintf("idle worker is not enough, need %d, but only %d", nums, m.count))
	}

	m.Lock()
	defer m.Unlock()

	if m.count < nums {
		return nil, errors.New(fmt.Sprintf("idle worker is not enough, need %d, but only %d", nums, m.count))
	}

	res := make([]*Worker, 0, nums)
	for i := 0; i < nums; i++ {
		res = append(res, m.PopWithoutLock())
	}

	return res, nil
}
