package pool

import (
	"awesomeProject/4_pool/worker"
	"fmt"
	"sync/atomic"
)

type Pool struct {
	workers      []*worker.Worker
	taskQueue    chan func()
	workersCount uint64
	queueSize    uint64
	queueLength  *atomic.Int64
	isClosed     atomic.Bool
	isPaused     atomic.Bool
}

type Pooler interface {
	Submit(act func(...interface{}), args ...interface{}) bool
	SubmitWait() bool
	Stop() bool
	StopWait() bool
	Pause() bool
	Resume() bool
}

func NewPool(workersCount, queueSize uint64) (*Pool, error) {
	if workersCount == 0 {
		return nil, fmt.Errorf("workersCount must be greater than 0")
	}

	pool := &Pool{
		taskQueue:    make(chan func(), queueSize),
		workersCount: workersCount,
		queueSize:    queueSize,
		queueLength:  &atomic.Int64{},
		isClosed:     atomic.Bool{},
		isPaused:     atomic.Bool{},
	}

	var index uint64 = 0
	for index = 0; index < workersCount; index++ {
		_worker := worker.NewWorker(index, pool.taskQueue, pool.queueLength)
		pool.workers = append(pool.workers, _worker)
		_worker.Start()
	}

	return pool, nil
}

func (pool *Pool) Submit(act func(...interface{}), args ...interface{}) bool {
	taskWithArgs := func() {
		act(args...)
	}

	if pool.queueLength.Load() == int64(pool.queueSize) || pool.isClosed.Load() {
		return false
	}
	pool.queueLength.Add(1)
	pool.taskQueue <- taskWithArgs
	return true
}

func (pool *Pool) SubmitWait() bool {
	if pool.isPaused.Load() {
		return false
	}
	for {
		if pool.queueLength.Load() == 0 {
			return true
		}
	}
}

func (pool *Pool) Stop() bool {
	if pool.isClosed.Load() {
		return false
	}
	pool.isClosed.Store(true)
	close(pool.taskQueue)
	for _, _worker := range pool.workers {
		_worker.Stop()
	}
	return true
}

func (pool *Pool) StopWait() bool {
	if pool.isClosed.Load() || pool.isPaused.Load() {
		return false
	}

	pool.isClosed.Store(true)
	close(pool.taskQueue)
	for {
		if pool.queueLength.Load() == 0 {
			return false
		}
	}
}

func (pool *Pool) Pause() bool {
	if pool.isPaused.Load() {
		return false
	}
	for _, _worker := range pool.workers {
		_worker.Pause()
	}
	pool.isPaused.Store(true)
	return true
}

func (pool *Pool) Resume() bool {
	if !pool.isPaused.Load() {
		return false
	}
	for _, _worker := range pool.workers {
		_worker.Resume()
	}
	pool.isPaused.Store(false)
	return true
}
