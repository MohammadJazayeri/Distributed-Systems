# Part 1 - IPC with Two Separate Processes

## Overview
This part implements a simple request/response service using two independent Go processes:
- interface.go: receives user input, sends requests to the worker, and displays responses
- worker.go: continuously listens for requests, processes them, and sends back results

Communication between the two processes is done using a Pipe.

## File Structure
- interface.go
- worker.go
- README.md

## Dependencies
- Go
- Linux
- Go standard library only

## Build
```
go build -o interface interface.go
go build -o worker worker.go
```

## Run
Run worker first, then interface:
```
./worker
./interface
```

## Behavior
1. User enters a request in the interface.
2. Interface sends the request through the pipe.
3. Worker processes the request.
4. Worker sends the result back.
5. Interface prints the result.
