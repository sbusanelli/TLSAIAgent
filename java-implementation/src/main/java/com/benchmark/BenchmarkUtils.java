package com.benchmark;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryMXBean;

public class BenchmarkUtils {
    
    public static class BenchmarkResult {
        public final String implementation;
        public final int taskCount;
        public final long executionTimeMs;
        public final long memoryUsedMB;
        public final double throughput;
        
        public BenchmarkResult(String implementation, int taskCount, long executionTimeMs, long memoryUsedMB) {
            this.implementation = implementation;
            this.taskCount = taskCount;
            this.executionTimeMs = executionTimeMs;
            this.memoryUsedMB = memoryUsedMB;
            this.throughput = (double) taskCount / (executionTimeMs / 1000.0);
        }
        
        @Override
        public String toString() {
            return String.format("%s: %d tasks in %dms (%.2f tasks/sec, %dMB memory)", 
                implementation, taskCount, executionTimeMs, throughput, memoryUsedMB);
        }
    }
    
    public static long getMemoryUsage() {
        MemoryMXBean memoryBean = ManagementFactory.getMemoryMXBean();
        return memoryBean.getHeapMemoryUsage().getUsed() / (1024 * 1024);
    }
    
    public static BenchmarkResult runBenchmark(String implementation, int taskCount, boolean useVirtualThreads) {
        System.gc(); // Clean up before benchmark
        long startMemory = getMemoryUsage();
        long startTime = System.nanoTime();
        
        try (ExecutorService executor = useVirtualThreads ? 
                Executors.newVirtualThreadPerTaskExecutor() : 
                Executors.newCachedThreadPool()) {
            
            List<CompletableFuture<Void>> futures = new ArrayList<>();
            
            for (int i = 0; i < taskCount; i++) {
                CompletableFuture<Void> future = CompletableFuture.runAsync(() -> {
                    simulateIOBoundTask();
                }, executor);
                futures.add(future);
            }
            
            // Wait for all tasks to complete
            CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
        }
        
        long endTime = System.nanoTime();
        long endMemory = getMemoryUsage();
        
        long executionTimeMs = (endTime - startTime) / 1_000_000;
        long memoryUsed = Math.max(0, endMemory - startMemory);
        
        return new BenchmarkResult(implementation, taskCount, executionTimeMs, memoryUsed);
    }
    
    private static void simulateIOBoundTask() {
        try {
            // Simulate I/O operation (network request, file read, etc.)
            Thread.sleep(50); // 50ms simulated I/O latency
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
    
    public static void printResults(List<BenchmarkResult> results) {
        System.out.println("\n=== BENCHMARK RESULTS ===");
        System.out.println(String.format("%-20s %-10s %-12s %-15s %-15s", 
            "Implementation", "Tasks", "Time (ms)", "Throughput", "Memory (MB)"));
        System.out.println("-".repeat(80));
        
        for (BenchmarkResult result : results) {
            System.out.println(String.format("%-20s %-10d %-12d %-15.2f %-15d", 
                result.implementation, result.taskCount, result.executionTimeMs, 
                result.throughput, result.memoryUsedMB));
        }
        System.out.println();
    }
    
    public static void runComparison(int taskCount) {
        System.out.printf("Running benchmark with %d I/O-bound tasks...\n\n", taskCount);
        
        List<BenchmarkResult> results = new ArrayList<>();
        
        // Run traditional threads benchmark
        BenchmarkResult traditional = runBenchmark("Traditional Threads", taskCount, false);
        results.add(traditional);
        
        // Small delay between tests
        try {
            Thread.sleep(1000);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        
        // Run virtual threads benchmark
        BenchmarkResult virtual = runBenchmark("Virtual Threads", taskCount, true);
        results.add(virtual);
        
        printResults(results);
        
        // Calculate performance improvement
        double improvement = (double) traditional.executionTimeMs / virtual.executionTimeMs;
        System.out.printf("Virtual Threads are %.2fx faster than Traditional Threads\n", improvement);
    }
}
