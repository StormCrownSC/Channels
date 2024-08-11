package main

import (
	"awesomeProject/4_pool/worker_pool"
	"fmt"
	"time"
)

func main() {
	// Create a pool with 4 worker threads and a queue size of 10.
	pool, err := worker_pool.NewPool(4, 10)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i := 0; i < 100; i++ {
		pool.Submit(func(args ...interface{}) {
			fmt.Println(fmt.Sprintf("Task executed with args %d", args))
			time.Sleep(time.Duration(i) * time.Second)
		}, i)
	}

	// Waiting for tasks to complete and shutting down the pool
	pool.Shutdown()

}
