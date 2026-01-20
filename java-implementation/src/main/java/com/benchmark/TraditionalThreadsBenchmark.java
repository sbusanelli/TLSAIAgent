package com.benchmark;

public class TraditionalThreadsBenchmark {
    public static void main(String[] args) {
        System.out.println("=== Java Traditional Threads Benchmark ===");
        System.out.println("Testing Traditional Threads performance with I/O-bound tasks\n");

        int[] taskCounts = {1000, 5000, 10000, 25000}; // Lower max due to OS thread limits

        for (int taskCount : taskCounts) {
            System.out.printf("Running %d tasks with Traditional Threads...\n", taskCount);
            BenchmarkUtils.BenchmarkResult result = BenchmarkUtils.runBenchmark("Traditional Threads", taskCount, false);
            System.out.println(result);
            System.out.println();

            // Longer delay for traditional threads to allow cleanup
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        System.out.println("Traditional Threads benchmark completed!");
    }
}
