package worker

import (
	"sync/atomic"
)

type Worker struct {
	id          uint64
	taskCh      chan func()
	closeCh     chan struct{}
	pauseCh     chan struct{}
	queueLength *atomic.Int64
}

type WorkerRunner interface {
	Start()
	Stop()
	Pause()
	Resume()
}

func NewWorker(id uint64, taskCh chan func(), queueLength *atomic.Int64) *Worker {
	return &Worker{
		id:          id,
		closeCh:     make(chan struct{}),
		pauseCh:     make(chan struct{}),
		taskCh:      taskCh,
		queueLength: queueLength,
	}
}

func (worker *Worker) Start() {
	go func() {
		for {
			select {
			case act := <-worker.taskCh:
				if act != nil {
					act()
				}
				worker.queueLength.Add(-1)
			case <-worker.pauseCh:
				select {
				case <-worker.pauseCh:
					continue
				case <-worker.closeCh:
					close(worker.pauseCh)
					close(worker.closeCh)
					return
				}
			case <-worker.closeCh:
				close(worker.pauseCh)
				close(worker.closeCh)
				return
			}
		}
	}()
}

func (worker *Worker) Stop() {
	worker.closeCh <- struct{}{}
}

func (worker *Worker) Pause() {
	worker.pauseCh <- struct{}{}
}

func (worker *Worker) Resume() {
	worker.pauseCh <- struct{}{}
}
