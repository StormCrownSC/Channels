package worker

import (
	"sync"
	"sync/atomic"
)

type Task struct {
	Action func(...interface{})
	Args   []interface{}
	WG     *sync.WaitGroup
}

type Worker struct {
	id          uint64
	taskCh      chan *Task
	closeCh     chan struct{}
	queueLength *atomic.Int64
}

type WorkerRunner interface {
	Start()
	Stop()
}

func NewWorker(id uint64, taskCh chan *Task, queueLength *atomic.Int64) *Worker {
	return &Worker{
		id:          id,
		closeCh:     make(chan struct{}),
		taskCh:      taskCh,
		queueLength: queueLength,
	}
}

func (worker *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-worker.taskCh:
				if task != nil {
					task.Action(task.Args...)
					if task.WG != nil {
						task.WG.Done()
					}
					worker.queueLength.Add(-1)
				}
			case <-worker.closeCh:
				return
			}
		}
	}()
}

func (worker *Worker) Stop() {
	worker.closeCh <- struct{}{}
}
