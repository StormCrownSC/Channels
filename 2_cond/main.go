package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

type Task struct {
	TaskName string
	TaskFunc func() error
}

type Queue struct {
	tasks chan Task
}

func NewQueue() *Queue {
	q := Queue{tasks: make(chan Task, runtime.NumCPU())}
	return &q
}

func (q *Queue) Push(t *Task) { // одобавить таску в очередь
	q.tasks <- *t
}

func (q *Queue) PopWait() *Task { // достаем таску из очереди, если тасок нет, блокируемся.
	t := <-q.tasks

	return &t
}

type Worker struct {
	id    int
	queue *Queue
}

func NewWorker(id int, queue *Queue) *Worker {
	w := Worker{
		id:    id,
		queue: queue,
	}
	return &w
}

func (w *Worker) Loop() {
	for {
		t := w.queue.PopWait()

		err := t.TaskFunc()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Printf("worker #%d resized %s\n", w.id, t.TaskName)
	}
}

func main() {
	queue := NewQueue()
	workers := make([]*Worker, 0, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		workers = append(workers, NewWorker(i, queue))
	}

	for _, w := range workers {
		go w.Loop()
	}
	tasks := []*Task{
		{
			TaskName: "fmt HelloWorld",
			TaskFunc: func() error {
				fmt.Println("Hello world")
				return nil
			},
		},
		{
			TaskName: "something with err",
			TaskFunc: func() error {
				return errors.New("test error")
			},
		},
	}

	queue.Push(tasks[0])
	queue.Push(tasks[1])

	time.Sleep(1 * time.Second)
}
