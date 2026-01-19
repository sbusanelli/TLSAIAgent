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
			taskCounts := []int{1000, 5000, 10000, 50000, 100000}
			
			if c.Int("tasks") > 0 {
				taskCounts = []int{c.Int("tasks")}
			}

			fmt.Println("🚀 Kicking off the Go Goroutine Benchmark...")
			fmt.Println("------------------------------------------")
			fmt.Printf("⚙️  Configuration: %d tasks, %s duration, %d workers\n\n",
				c.Int("tasks"), c.String("duration"), c.Int("workers"))

			// Run Go goroutine benchmarks
			runGoroutineComparison(taskCounts)

			fmt.Println("\n✅ Go Goroutine benchmark complete!")
			fmt.Println("-----------------------------------")
			fmt.Println("\nNext, compare with Java Virtual Threads:")
			fmt.Println("  cd ../java-implementation && mvn exec:java -Dexec.mainClass=\"com.benchmark.ComparisonRunner\"")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
