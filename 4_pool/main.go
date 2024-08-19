package main

import (
	"awesomeProject/4_pool/pool"
	"context"
	"fmt"
	"math/rand"
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
	for i := 0; i < 100; i++ {
		success := _pool.Submit(func(args ...interface{}) {
			fmt.Println(fmt.Sprintf("Task executed with args %d", args))
			time.Sleep(time.Duration(rand.Intn(i+1)*100) * time.Millisecond)
		}, i)
		if !success {
			fmt.Println("task executed failed")
		}
		//time.Sleep(time.Duration(rand.Intn(i+1)*10) * time.Millisecond)
	}

	fmt.Println("Pool pause")
	_pool.Pause()
	time.Sleep(1 * time.Second)
	_pool.Resume()
	fmt.Println("Pool resume")

	// Waiting for tasks to complete and shutting down the pool
	_pool.StopWait()
	stopChan <- struct{}{}
}
