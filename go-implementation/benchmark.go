package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
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
	fmt.Println("\n📊 Go Goroutine Benchmark Results")
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("%-15s | %-10s | %-15s | %-15s | %-15s\n", "Implementation", "Tasks", "Time", "Memory (MB)", "Throughput (t/s)")
	fmt.Println("-----------------------------------------------------------------------------")

	for _, result := range results {
		throughputStr := fmt.Sprintf("%.2f", result.Throughput)
		fmt.Printf("%-15s | %-10d | %-15v | %-15d | %-15s\n",
			result.Implementation, result.TaskCount, result.ExecutionTime, result.MemoryUsed, throughputStr)
	}
	fmt.Println("-----------------------------------------------------------------------------")
}

func printScalabilityAnalysis(results []BenchmarkResult) {
	fmt.Println("\n📈 Goroutine Scalability Analysis")
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("%-20s | %-20s | %-20s\n", "Task Increase", "Execution Time Increase", "Tasks-to-Time Ratio")
	fmt.Println("----------------------------------------------------------------")

	for i, result := range results {
		if i > 0 {
			prevResult := results[i-1]

			if prevResult.TaskCount == 0 || prevResult.ExecutionTime.Seconds() == 0 {
				continue
			}

			taskIncrease := float64(result.TaskCount) / float64(prevResult.TaskCount)
			timeIncrease := result.ExecutionTime.Seconds() / prevResult.ExecutionTime.Seconds()

			ratio := 0.0
			if timeIncrease > 0 {
				ratio = taskIncrease / timeIncrease
			}

			taskIncreaseStr := fmt.Sprintf("%.2fx", taskIncrease)
			timeIncreaseStr := fmt.Sprintf("%.2fx", timeIncrease)
			ratioStr := fmt.Sprintf("%.2f", ratio)

			fmt.Printf("%-20s | %-20s | %-20s\n", taskIncreaseStr, timeIncreaseStr, ratioStr)
		}
	}
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("* A Tasks-to-Time Ratio > 1.0 indicates good scalability.")
}

func runGoroutineComparison(taskCounts []int) {
	fmt.Println("🏃 Running Go Goroutine benchmarks...")
	var results []BenchmarkResult

	for _, taskCount := range taskCounts {
		// Simple progress indicator
		fmt.Printf("  • Running test with %d tasks...", taskCount)
		result := runGoroutineBenchmark(taskCount)
		results = append(results, result)
		fmt.Println(" Done.")

		// Brief pause between tests
		time.Sleep(250 * time.Millisecond)
	}

	printResults(results)

	if len(results) > 1 {
		printScalabilityAnalysis(results)
	}
}
