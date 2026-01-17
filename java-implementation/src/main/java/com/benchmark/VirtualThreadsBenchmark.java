package com.benchmark;

public class VirtualThreadsBenchmark {
    public static void main(String[] args) {
        System.out.println("=== Java Virtual Threads Benchmark ===");
        System.out.println("Testing Virtual Threads performance with I/O-bound tasks\n");
        
        int[] taskCounts = {1000, 5000, 10000, 50000};
        
        for (int taskCount : taskCounts) {
            BenchmarkUtils.runComparison(taskCount);
            System.out.println();
        }
        
        System.out.println("Virtual Threads benchmark completed!");
    }
}
