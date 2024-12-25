package messaging

import (
	"fmt"
)

type WorkerPool struct {
	workers   int
	jobQueue  chan []byte
	messenger Messenger
}

// Constructor
func NewWorkerPool(workers int, messenger Messenger) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		jobQueue:  make(chan []byte, 1000), // Tamaño de búfer ajustable
		messenger: messenger,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	fmt.Printf("Worker %d started\n", id)
	for job := range wp.jobQueue {
		err := wp.messenger.Send(job)
		if err != nil {
			fmt.Printf("Worker %d failed to send message: %v\n", id, err)
		}
	}
}

func (wp *WorkerPool) AddJob(data []byte) {
	wp.jobQueue <- data
}
