package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "virtual-threads-benchmark",
		Usage: "Performance comparison between Go Goroutines and Java Virtual Threads",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "tasks",
				Value: 10000,
				Usage: "Number of concurrent tasks to run",
			},
			&cli.StringFlag{
				Name:  "duration",
				Value: "50ms",
				Usage: "Duration of each simulated I/O task",
			},
			&cli.IntFlag{
				Name:  "workers",
				Value: 0, // 0 means use runtime.NumCPU() * 2
				Usage: "Maximum number of concurrent workers",
			},
		},
		Action: func(c *cli.Context) error {
			taskCounts := []int{1000, 5000, 10000} // Default tasks for comparison

			// If the user specifies a task count, override the default list
			if c.IsSet("tasks") {
				taskCounts = []int{c.Int("tasks")}
			}

			fmt.Println("🚀 Virtual Threads vs Goroutines Benchmark")
			fmt.Println("=====================================")
			fmt.Printf("Configuration: %d tasks max, %s duration per task\n\n",
				c.Int("tasks"), c.String("duration"))

			// Run Go goroutine benchmarks
			runGoroutineComparison(taskCounts)

			fmt.Println("\n📊 Benchmark completed!")
			fmt.Println("Run the Java implementation to compare results:")
			fmt.Println("  cd ../java-implementation")
			fmt.Println("  mvn exec:java -Dexec.mainClass=\"com.benchmark.ComparisonRunner\"")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
