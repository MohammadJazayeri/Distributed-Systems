# Automated Test Client

This directory contains the automated test benchmarking client designed to evaluate the performance, resilience, and convergence of the distributed key-value store under different network conditions.

## Automated Scenarios
The client executes 4 distinct testing pipelines sequentially:
1. **Scenario 1 (Temporary Inconsistency):** Writes data eventually with a high network delay and immediately performs a read from another replica to capture stale reads and calculate system convergence time.
2. **Scenario 2 (Replica Failure Analysis):** Measures the impact of node failures. If you manually kill a replica (e.g., node 3), it demonstrates how strong consistency failing quorum blocks writes while eventual consistency remains highly available.
3. **Scenario 3 (Concurrent Conflicts):** Dispatches concurrent conflicting values to different replicas simultaneously to verify that the Last-Write-Wins (LWW) resolution algorithm correctly converges all nodes to a single consistent state.
4. **Scenario 4 (Network Latency Benchmarks):** Evaluates a comprehensive matrix comparing Strong vs Eventual consistency models under varying injected delays (0ms, 500ms, 2000ms).

## Tracked Metrics
The client captures and evaluates the following benchmarking metrics:
- **PUT Latency:** Wall-clock duration of write requests.
- **System Convergence Time:** Total time taken for all replicas to catch up to the latest state.
- **Stale Reads Count:** Total number of out-of-date read operations captured during a synchronization window.
- **Immediate Replicas Updated:** Total number of replicas synchronized prior to client HTTP acknowledgement.

## How to Run
Once all 3 replica nodes are up and running in their respective terminals, open a new terminal in this directory and execute:

```bash
go run main.go
Outputs and Logging
Upon successful execution, a results/ folder will be generated in the root project directory containing raw metrics logs:

scenario1.txt: Inconsistency capturing window logs.

scenario2.txt: Quorum stability and failure evaluation.

scenario3.txt: Final resolved value state.

scenario4.txt: Complete benchmarking matrix text logs, suitable for populating the evaluation reports.