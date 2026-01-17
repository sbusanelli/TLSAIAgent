package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"os"
	"strconv"
)

type BenchmarkResult struct {
	Implementation string
	TaskCount     int
	ExecutionTime time.Duration
	MemoryUsed   int64
	Throughput   float64
}

func (br BenchmarkResult) String() string {
	return fmt.Sprintf("%s: %d tasks in %v (%.2f tasks/sec, %d MB memory)", 
		br.Implementation, br.TaskCount, br.ExecutionTime, br.Throughput, br.MemoryUsed)
}

func getMemoryUsage() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Alloc) / 1024 / 1024 // Convert to MB
}

func simulateIOBoundTask() {
	time.Sleep(50 * time.Millisecond) // 50ms simulated I/O latency
}

func runGoroutineBenchmark(taskCount int) BenchmarkResult {
	runtime.GC() // Clean up before benchmark
	startMemory := getMemoryUsage()
	startTime := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.NumCPU()*2) // Limit concurrent goroutines

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			simulateIOBoundTask()
		}()
	}

	wg.Wait()

	endTime := time.Now()
	endMemory := getMemoryUsage()

	executionTime := endTime.Sub(startTime)
	memoryUsed := endMemory - startMemory
	if memoryUsed < 0 {
		memoryUsed = 0
	}

	throughput := float64(taskCount) / executionTime.Seconds()

	return BenchmarkResult{
		Implementation: "Go Goroutines",
		TaskCount:     taskCount,
		ExecutionTime: executionTime,
		MemoryUsed:   memoryUsed,
		Throughput:   throughput,
	}
}

func printResults(results []BenchmarkResult) {
	fmt.Println("\n=== GOROUTINE BENCHMARK RESULTS ===")
	fmt.Printf("%-20s %-10s %-15s %-15s\n", "Implementation", "Tasks", "Time", "Throughput")
	fmt.Println("-".repeat(60))

	for _, result := range results {
		fmt.Printf("%-20s %-10d %-15v %-15.2f\n", 
			result.Implementation, result.TaskCount, result.ExecutionTime, result.Throughput)
	}
	fmt.Println()
}

func runGoroutineComparison(taskCounts []int) {
	fmt.Println("=== Go Goroutines Performance Test ===")
	fmt.Println("Testing Goroutines performance with I/O-bound tasks\n")

	var results []BenchmarkResult

	for _, taskCount := range taskCounts {
		fmt.Printf("Running %d tasks with Goroutines...\n", taskCount)
		result := runGoroutineBenchmark(taskCount)
		results = append(results, result)
		fmt.Println(result)
		fmt.Println()

		// Brief pause between tests
		time.Sleep(500 * time.Millisecond)
	}

	printResults(results)

	fmt.Println("=== Goroutine Scalability Analysis ===")
	for i, result := range results {
		if i > 0 {
			prevResult := results[i-1]
			scalability := float64(result.TaskCount) / float64(prevResult.TaskCount)
			timeIncrease := result.ExecutionTime.Seconds() / prevResult.ExecutionTime.Seconds()
			
			fmt.Printf("Task count increased %.1fx, execution time increased %.2fx\n", 
				scalability, timeIncrease)
		}
	}
}
