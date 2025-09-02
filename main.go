package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Task represents work to be done
type Task struct {
	ID   int
	Data string
}

// Result represents the outcome of processing a task
type Result struct {
	TaskID   int
	Output   string
	Error    error
	WorkerID int
}

// ProcessTask simulates computational work (DO NOT MODIFY)
func ProcessTask(ctx context.Context, task Task, wg *sync.WaitGroup) (string, error) {
	defer wg.Done()
	fmt.Println("starting task ", task.ID)

	// Simulate variable processing time (50-200ms)
	processingTime := time.Duration(50+rand.Intn(150)) * time.Millisecond

	select {
	case <-time.After(processingTime):
		// Simulate 10% failure rate
		if rand.Intn(10) == 0 {
			return "", fmt.Errorf("processing failed for task %d", task.ID)
		}
		// Simulate some computation
		output := fmt.Sprintf("processed-%s-result-%d", task.Data, task.ID*2)
		return output, nil
	case <-ctx.Done():
		return "", fmt.Errorf("for some reason, task %d  was canceled: %v", task.ID, ctx.Err())
	}
}

// WorkerPool manages the pool of workers
type WorkerPool struct {
	Count int
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{Count: numWorkers}
}

// ProcessTasks processes all tasks using the worker pool
// Returns results in the same order as input tasks
func (wp *WorkerPool) ProcessTasks(ctx context.Context, tasks []Task) ([]Result, error) {
	taskChan := make(chan Task, len(tasks))
	resultChan := make(chan Result, len(tasks))

	// starting workers
	var wg sync.WaitGroup
	for i := 0; i < wp.Count; i++ {
		wg.Add(1)
		go wp.worker(ctx, i, taskChan, resultChan, &wg)
	}

	// sending tasks to workers (yeah, I had to use channel)
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// waiting for everything to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// collecting the results on channel
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}
	sortResults(results)

	return results, nil
}

func sortResults(results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].TaskID < results[j].TaskID
	})
}

// worker processes tasks from the task channel
func (wp *WorkerPool) worker(ctx context.Context, workerID int, taskChan <-chan Task, resultChan chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		var taskWg sync.WaitGroup
		taskWg.Add(1)
		output, err := ProcessTask(ctx, task, &taskWg)

		resultChan <- Result{
			TaskID:   task.ID,
			Output:   output,
			Error:    err,
			WorkerID: workerID,
		}
	}
}

func main() {
	// Seed random number generator for consistent testing
	rand.Seed(time.Now().UnixNano())

	// Create test tasks
	tasks := make([]Task, 20)
	for i := 0; i < 20; i++ {
		tasks[i] = Task{
			ID:   i,
			Data: fmt.Sprintf("data-%d", i),
		}
	}

	// Create worker pool with 3 workers
	pool := NewWorkerPool(3)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Processing %d tasks with 3 workers...\n", len(tasks))
	start := time.Now()

	results, err := pool.ProcessTasks(ctx, tasks)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Analyze results
	successful := 0
	failed := 0

	fmt.Printf("\n=== Results ===\n")
	for _, result := range results {
		if result.Error == nil {
			successful++
			fmt.Printf("Task %d: SUCCESS (Worker %d) -> %s\n",
				result.TaskID, result.WorkerID, result.Output)
		} else {
			failed++
			fmt.Printf("Task %d: FAILED (Worker %d) -> %v\n",
				result.TaskID, result.WorkerID, result.Error)
		}
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total tasks: %d\n", len(results))
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Processing time: %v\n", duration)
	fmt.Printf("Success rate: %.1f%%\n", float64(successful)/float64(len(results))*100)
}
