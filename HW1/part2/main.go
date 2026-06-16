package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

type WorkloadType string

const (
	CPU_BOUND WorkloadType = "CPU_BOUND"
	MIXED     WorkloadType = "MIXED"
)

type ExperimentConfig struct {
	GoroutineCount int
	Workload       WorkloadType
	GOMAXPROCS     int
	TraceFile      string
}

type ExperimentResult struct {
	GoroutineCount      int
	Workload            WorkloadType
	GOMAXPROCS          int
	ExecutionTime       time.Duration
	Throughput          float64
	AvgLatency          time.Duration
	MinLatency          time.Duration
	MaxLatency          time.Duration
	StdDevLatency       time.Duration
	GoroutineCountSeen  int
}

type TaskTiming struct {
	Start time.Time
	End   time.Time
}

func cpuBoundWork(iterations int) {
	var sum float64
	for i := 0; i < iterations; i++ {
		sum += math.Sqrt(float64(i)) * 1.000001
	}
	_ = sum
}

func mixedWork(iterations int) {
	var sum float64
	for i := 0; i < iterations; i++ {
		sum += math.Sqrt(float64(i)) * 1.000001
	}
	_ = sum

	time.Sleep(5 * time.Millisecond)
}

func runExperiment(config ExperimentConfig) ExperimentResult {
	log.Printf("Starting experiment: goroutines=%d workload=%s GOMAXPROCS=%d",
		config.GoroutineCount, config.Workload, config.GOMAXPROCS)

	oldProcs := runtime.GOMAXPROCS(config.GOMAXPROCS)
	_ = oldProcs

	var traceFile *os.File
	if config.TraceFile != "" {
		f, err := os.Create(config.TraceFile)
		if err != nil {
			log.Fatalf("failed to create trace file: %v", err)
		}
		traceFile = f
		if err := trace.Start(traceFile); err != nil {
			_ = traceFile.Close()
			log.Fatalf("failed to start trace: %v", err)
		}
	}

	if traceFile != nil {
		defer func() {
			trace.Stop()
			_ = traceFile.Close()
		}()
	}

	timings := make([]TaskTiming, config.GoroutineCount)
	var wg sync.WaitGroup

	startAll := time.Now()

	for i := 0; i < config.GoroutineCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			taskStart := time.Now()

			switch config.Workload {
			case CPU_BOUND:
				cpuBoundWork(15_000_000)
			case MIXED:
				mixedWork(2_000_000)
			default:
				log.Fatalf("unknown workload: %s", config.Workload)
			}

			taskEnd := time.Now()
			timings[id] = TaskTiming{Start: taskStart, End: taskEnd}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(startAll)

	latencies := make([]time.Duration, 0, config.GoroutineCount)
	var minLatency, maxLatency time.Duration

	for i, t := range timings {
		_ = i
		if t.End.IsZero() || t.Start.IsZero() {
			continue
		}
		lat := t.End.Sub(t.Start)
		latencies = append(latencies, lat)

		if len(latencies) == 1 || lat < minLatency {
			minLatency = lat
		}
		if len(latencies) == 1 || lat > maxLatency {
			maxLatency = lat
		}
	}

	var sum time.Duration
	for _, lat := range latencies {
		sum += lat
	}

	avgLatency := time.Duration(0)
	if len(latencies) > 0 {
		avgLatency = sum / time.Duration(len(latencies))
	}

	var variance float64
	if len(latencies) > 1 {
		mean := float64(avgLatency)
		for _, lat := range latencies {
			diff := float64(lat) - mean
			variance += diff * diff
		}
		variance /= float64(len(latencies) - 1)
	}

	stdDev := time.Duration(math.Sqrt(variance))

	throughput := 0.0
	if totalTime.Seconds() > 0 {
		throughput = float64(config.GoroutineCount) / totalTime.Seconds()
	}

	result := ExperimentResult{
		GoroutineCount:     config.GoroutineCount,
		Workload:           config.Workload,
		GOMAXPROCS:         config.GOMAXPROCS,
		ExecutionTime:      totalTime,
		Throughput:         throughput,
		AvgLatency:         avgLatency,
		MinLatency:         minLatency,
		MaxLatency:         maxLatency,
		StdDevLatency:      stdDev,
		GoroutineCountSeen: runtime.NumGoroutine(),
	}

	log.Printf("Finished: time=%v throughput=%.2f avg=%v min=%v max=%v stddev=%v",
		result.ExecutionTime, result.Throughput, result.AvgLatency, result.MinLatency, result.MaxLatency, result.StdDevLatency)

	return result
}

func main() {
	initial := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	log.Printf("Initial GOMAXPROCS=%d NumCPU=%d", initial, numCPU)

	gCounts := []int{1, 2, 4, 8, 16, 32, 64}
	procs := []int{1, 2, numCPU}
	workloads := []WorkloadType{CPU_BOUND, MIXED}

	results := make([]ExperimentResult, 0)

	for _, w := range workloads {
		for _, p := range procs {
			for _, g := range gCounts {
				cfg := ExperimentConfig{
					GoroutineCount: g,
					Workload:       w,
					GOMAXPROCS:     p,
				}

				// Save a few traces only, to avoid huge disk usage.
				if g == 1 && (p == 1 || p == 2 || p == numCPU) {
					cfg.TraceFile = fmt.Sprintf("trace_%s_g%d_p%d.out", w, g, p)
				}

				results = append(results, runExperiment(cfg))
			}
		}
	}

	fmt.Println()
	fmt.Println("SUMMARY")
	fmt.Println("Workload\tGOMAXPROCS\tGoroutines\tTotalTime\tThroughput\tAvgLatency\tMinLatency\tMaxLatency\tStdDevLatency")
	for _, r := range results {
		fmt.Printf("%s\t%d\t%d\t%v\t%.2f\t%v\t%v\t%v\t%v\n",
			r.Workload,
			r.GOMAXPROCS,
			r.GoroutineCount,
			r.ExecutionTime,
			r.Throughput,
			r.AvgLatency,
			r.MinLatency,
			r.MaxLatency,
			r.StdDevLatency,
		)
	}
}
