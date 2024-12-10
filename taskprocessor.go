package taskprocessor

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work to be processed.
type Task func() error

// Processor manages a pool of workers to process tasks concurrently.
type Processor struct {
	workers      int
	taskQueue    chan Task
	wg           sync.WaitGroup
	shutdownChan chan struct{}
	rateLimiter  *time.Ticker
}

// NewProcessor creates a new task processor with the specified number of workers.
func NewProcessor(workers int, rateLimit time.Duration) *Processor {
	p := &Processor{
		workers:      workers,
		taskQueue:    make(chan Task, workers*10), // Buffer size to prevent blocking
		shutdownChan: make(chan struct{}),
	}

	if rateLimit > 0 {
		p.rateLimiter = time.NewTicker(rateLimit)
	}

	// Start worker goroutines
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

// Submit adds a task to the queue for processing.
func (p *Processor) Submit(task Task) {
	p.taskQueue <- task
}

// Stop shuts down the processor and waits for all tasks to complete.
func (p *Processor) Stop() {
	close(p.shutdownChan)
	p.wg.Wait()
}

func (p *Processor) worker() {
	defer p.wg.Done()

	for {
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				return // Shutdown signal received
			}

			if p.rateLimiter != nil {
				<-p.rateLimiter.C
			}

			if err := task(); err != nil {
				// Log the error
				fmt.Printf("Error processing task: %v\n", err)
			}

		case <-p.shutdownChan:
			return // Shutdown signal received
		}
	}
}
