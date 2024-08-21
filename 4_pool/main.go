package main

import (
	"awesomeProject/4_pool/pool"
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a pool with 4 worker threads and a queue size of 10.
	_pool, err := pool.NewPool(4, 100)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var stopChan = make(chan struct{})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	go poolExecutor(stopChan, _pool)

	select {
	case <-ctx.Done():
		fmt.Println("stopping server...")
		_pool.Stop()
		fmt.Println("server is stopped")
	case <-stopChan:
		fmt.Println("stopping server...")
		fmt.Println("server is stopped")
	}
}

func poolExecutor(stopChan chan struct{}, _pool *pool.Pool) {
	fmt.Println("Start waiting task")
	_pool.SubmitWait(func(args ...interface{}) {
		fmt.Println(fmt.Sprintf("Task executed with args %d", args))
		time.Sleep(1 * time.Second)
	}, 1000)
	fmt.Println("Execute waiting task")

	for i := 0; i < 100; i++ {
		err := _pool.Submit(func(args ...interface{}) {
			fmt.Println(fmt.Sprintf("Task executed with args %d", args))
			time.Sleep(100 * time.Millisecond)
		}, i)
		if err != nil {
			fmt.Println(fmt.Sprintf("task executed failed [%s]", err))
		}
	}

	// Waiting for tasks to complete and shutting down the pool
	_pool.StopWait()
	stopChan <- struct{}{}
}
