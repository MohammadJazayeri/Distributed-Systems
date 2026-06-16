# Part 3 - HTTP Service in Docker

## Overview
This part implements a Go HTTP service, containerizes it with Docker, and runs it inside a virtual machine.

## File Structure
- main.go
- Dockerfile
- README.md

## Dependencies
- Go
- Docker
- Linux

## Endpoints

GET /health
Returns service health status.

GET /compute
Example request:
/compute?op=add&a=5&b=7

Example response:
{
  "operation": "add",
  "a": 5,
  "b": 7,
  "result": 12
}

## Error Handling
- Missing parameters
- Invalid operation
- Non‑numeric inputs
- Division by zero
- Wrong HTTP method

## Run without Docker
```
go run main.go
```

## Build Docker Image
```
docker build -t part3-service .
```

## Run Container
```
docker run -p 8080:8080 part3-service
```

Access from host:
http://localhost:8080/health
http://localhost:8080/compute?op=add&a=5&b=7
