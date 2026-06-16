# Part 2 - Concurrency and Scheduling Experiment

## Overview
This program evaluates the effect of increasing concurrency in Go using goroutines and different GOMAXPROCS settings.

## File Structure
- main.go
- README.md

## Dependencies
- Go
- Linux
- Go standard library only

## Experiment Parameters
Goroutine counts tested:
1, 2, 4, 8, 16, 32, 64

GOMAXPROCS values:
- 1
- 2
- runtime.NumCPU()

## Workloads
1. Computational workload (CPU-bound loop calculations)
2. Mixed workload (computation + sleep/channel/mutex)

## Metrics Collected
- Total execution time
- Throughput
- Average latency per task
- Number of goroutines
- GOMAXPROCS value

## Run
```
go run main.go
```

Or:
```
go build -o part2 main.go
./part2
```
