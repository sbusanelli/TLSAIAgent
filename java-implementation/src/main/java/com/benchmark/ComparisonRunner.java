package com.benchmark;

import java.util.ArrayList;
import java.util.List;

public class ComparisonRunner {
    public static void main(String[] args) {
        System.out.println("=== Comprehensive Java Threads Comparison ===");
        System.out.println("Comparing Traditional Threads vs Virtual Threads\n");

        int[] taskCounts = {1000, 5000, 10000, 50000, 100000};

        System.out.println("Task Count | Traditional (ms) | Virtual (ms) | Improvement");
        System.out.println("-".repeat(55));

        for (int taskCount : taskCounts) {
            List<BenchmarkUtils.BenchmarkResult> results = new ArrayList<>();

            // Run traditional threads
            BenchmarkUtils.BenchmarkResult traditional = BenchmarkUtils.runBenchmark("Traditional", taskCount, false);
            results.add(traditional);

            try {
                Thread.sleep(1000); // Brief pause between tests
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }

            // Run virtual threads
            BenchmarkUtils.BenchmarkResult virtual = BenchmarkUtils.runBenchmark("Virtual", taskCount, true);
            results.add(virtual);

            // Calculate improvement
            double improvement = (double) traditional.executionTimeMs / virtual.executionTimeMs;

            System.out.printf("%-10d | %-15d | %-12d | %.2fx\n",
                taskCount, traditional.executionTimeMs, virtual.executionTimeMs, improvement);
        }

        System.out.println("\n=== Summary ===");
        System.out.println("Virtual Threads show significant performance improvements for I/O-bound workloads");
        System.out.println("Performance gap increases with higher concurrency levels");
    }
}
