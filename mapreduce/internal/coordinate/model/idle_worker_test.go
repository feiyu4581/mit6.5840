package model

import (
	"fmt"
	"testing"
	"time"
)

func TestNewIdleWorkerManager(t *testing.T) {
	now := time.Now()
	manager := NewIdleWorkerManager()
	go func() {
		for i := 0; i < 20; i++ {
			manager.Add(&Worker{Status: WorkerRunningStatus})
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(10 * time.Minute)
	}()

	_ = manager.ChooseWorkers(5)
	fmt.Printf("CHoose 5 workers: %s\n", time.Now().Sub(now))

	_ = manager.ChooseWorkers(5)
	fmt.Printf("CHoose 5 workers: %s\n", time.Now().Sub(now))

	_ = manager.ChooseWorkers(5)
	fmt.Printf("CHoose 5 workers: %s\n", time.Now().Sub(now))

	_ = manager.ChooseWorkers(5)
	fmt.Printf("CHoose 5 workers: %s\n", time.Now().Sub(now))
}
