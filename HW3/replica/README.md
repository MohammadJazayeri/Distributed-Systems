# Replica Node - Distributed Key-Value Store

This directory contains the core server implementation for the individual replica nodes in the distributed key-value store system. Each node runs as an independent OS process and communicates with other replicas over HTTP.

## Features
- **In-Memory Storage:** Efficient key-value storage using Go's native `map` protected by `sync.RWMutex` for concurrent safety.
- **Consistency Models:**
  - **Strong Consistency:** Achieved via a synchronous Quorum write mechanism requiring acknowledgment from a majority of nodes (at least 2 out of 3).
  - **Eventual Consistency:** Performed asynchronously using Go Goroutines to return an immediate local acknowledgment while propagating updates in the background.
- **Versioning & Conflict Resolution:** Implements the Last-Write-Wins (LWW) strategy based on an incremental `Version` number and lexicographical node ID (`UpdatedBy`) comparison. This safely prevents stale network delayed packets from overwriting newer updates.

## Prerequisites
- [Go](https://go.dev/) (version 1.16 or higher)

## Configuration
Each replica node reads its network environment configuration from a JSON file located in the `configs/` folder. Example format:
```json
{
    "id": "replica1",
    "port": "8081",
    "peers": [
        "http://localhost:8082",
        "http://localhost:8083"
    ]
}
How to Run
To simulate a real distributed cluster, open 3 separate terminals and run each replica node with its respective configuration file:

Terminal 1 (Replica 1 - Port 8081):

Bash
go run main.go -config ../configs/replica1.json
Terminal 2 (Replica 2 - Port 8082):

Bash
go run main.go -config ../configs/replica2.json
Terminal 3 (Replica 3 - Port 8083):

Bash
go run main.go -config ../configs/replica3.json
API Endpoints
GET /get?key=<key>: Fetches the data item along with its version metadata.

POST /put: Coordinates a new write operation based on specified consistency level (strong/eventual) and network simulation delay (delay_ms).

POST /replicate: Internal cluster endpoint used for data propagation and conflict resolution execution.