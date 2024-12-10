package main

import (
	"rlhf/taskprocessor"
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	processor := taskprocessor.NewProcessor(10, 1*time.Second) // 10 workers, 1 task per second
	defer processor.Stop()

	for i := 0; i < 100; i++ {
		processor.Submit(func() error {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return nil
		})
	}

	time.Sleep(5 * time.Second) // Wait for tasks to complete
}
