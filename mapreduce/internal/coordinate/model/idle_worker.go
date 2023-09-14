package model

import (
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

	cond  *sync.Cond
	count int
}

func NewIdleWorkerManager() *IdleWorkerManager {
	manager := &IdleWorkerManager{}
	manager.cond = sync.NewCond(manager)
	return manager
}

func (m *IdleWorkerManager) Add(worker *Worker) {
	m.Lock()
	defer m.Unlock()
	defer m.cond.Broadcast()

	idleWorker := &IdleWorker{
		Worker: worker,
	}

	if m.IsEmpty() {
		m.Head = idleWorker
		m.Tail = idleWorker
		m.count = 1
		return
	}

	m.Tail.Next = idleWorker
	m.Tail = idleWorker
	m.count += 1
}

func (m *IdleWorkerManager) isCountEnough(count int) bool {
	if m.count < count {
		return false
	}

	currentCount := 0
	current := m.Head
	for current != nil {
		if current.Worker.IsRunningStatus() {
			currentCount += 1
			if currentCount >= count {
				return true
			}
		}

		current = current.Next
	}

	return false
}

func (m *IdleWorkerManager) PopActiveWorkerWithoutLock() *Worker {
	worker := m.PopWithoutLock()
	if worker == nil || worker.IsRunningStatus() {
		return worker
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

func (m *IdleWorkerManager) ChooseWorkers(nums int) []*Worker {
	m.Lock()
	defer m.Unlock()

	for !m.isCountEnough(nums) {
		m.cond.Wait()
	}

	res := make([]*Worker, 0, nums)
	for i := 0; i < nums; i++ {
		res = append(res, m.PopActiveWorkerWithoutLock())
	}

	return res
}
