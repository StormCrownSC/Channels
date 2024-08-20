package pool

import (
	"awesomeProject/4_pool/worker"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	workers      []*worker.Worker
	taskQueue    chan *worker.Task
	workersCount uint64
	queueSize    uint64
	queueLength  *atomic.Int64
	isClosed     atomic.Bool
}

type Pooler interface {
	Submit(act func(...any), args ...any) error
	SubmitWait(act func(...any), args ...any) error
	Stop() error
	StopWait() error
}

func NewPool(workersCount, queueSize uint64) (*Pool, error) {
	if workersCount == 0 {
		return nil, fmt.Errorf("workers count must be greater than 0")
	}

	pool := &Pool{
		taskQueue:    make(chan *worker.Task, queueSize),
		workersCount: workersCount,
		queueSize:    queueSize,
		queueLength:  &atomic.Int64{},
		isClosed:     atomic.Bool{},
	}

	var index uint64 = 0
	for index = 0; index < workersCount; index++ {
		_worker := worker.NewWorker(index, pool.taskQueue, pool.queueLength)
		pool.workers = append(pool.workers, _worker)
		_worker.Start()
	}

	return pool, nil
}

func (pool *Pool) Submit(act func(...interface{}), args ...interface{}) error {
	if pool.isClosed.Load() {
		return fmt.Errorf("worker pool is closed")
	}

	task := &worker.Task{
		Action: act,
		Args:   args,
		WG:     nil,
	}

	for {
		if pool.queueLength.Load() == int64(pool.queueSize) {
			time.Sleep(100 * time.Nanosecond)
			continue
		}
		pool.queueLength.Add(1)
		pool.taskQueue <- task
		break
	}

	return nil
}

func (pool *Pool) SubmitWait(act func(...interface{}), args ...interface{}) error {
	if pool.isClosed.Load() {
		return fmt.Errorf("worker pool is closed")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	task := &worker.Task{
		Action: act,
		Args:   args,
		WG:     &wg,
	}

	for {
		if pool.queueLength.Load() == int64(pool.queueSize) {
			time.Sleep(100 * time.Nanosecond)
			continue
		}
		pool.queueLength.Add(1)
		pool.taskQueue <- task
		break
	}

	wg.Wait() // Ожидание выполнения задачи

	return nil
}

func (pool *Pool) Stop() error {
	if pool.isClosed.Load() {
		return fmt.Errorf("worker pool already is closed")
	}
	pool.isClosed.Store(true)
	close(pool.taskQueue)
	for _, _worker := range pool.workers {
		_worker.Stop()
	}
	return nil
}

func (pool *Pool) StopWait() error {
	if pool.isClosed.Load() {
		return fmt.Errorf("worker pool already is closed")
	}

	pool.isClosed.Store(true)
	close(pool.taskQueue)
	for {
		if pool.queueLength.Load() == 0 {
			return nil
		}
		time.Sleep(100 * time.Nanosecond)
	}
}
