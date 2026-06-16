# Distributed Web Server

This project implements the web gateway and core orchestration layer of the distributed system. It handles user interactions, manages authentication via external RPCs, fetches protected assets, and monitors its own memory usage.

Processes:

1. Web Server Process (Publisher)

The Web Server serves a login portal, routes credentials to the authentication server, retrieves secure data upon success, and runs a continuous background telemetry loop to report high memory usage.

-------------------------------------

## Supported Operations / Endpoints

GET  /login          Displays the interactive HTML login portal.
POST /login          Authenticates user credentials and fetches a secure image if successful.
GET  /consume-memory Allocates 50MB of memory into the heap to simulate a memory leak.

-------------------------------------

## Response Format

Success (HTTP 200 OK):
Returns a stylized responsive HTML page rendering the user's secure desktop dashboard along with a Base64 encoded secure asset fetched from storage.

Errors:
HTTP 500: Critical Error: Auth Service (VM2) is offline.
HTML Render Error: Invalid credentials
HTML Render Error: VM3 is down for image fetch

-------------------------------------

## RPC & Network Ports Used

192.168.56.104:50051 (VM1 -> VM2 Auth Service via JSON-RPC)
192.168.56.102:50052 (VM1 -> VM3 File Service via Go RPC)
192.168.56.104:50053 (VM1 -> Dedicated Monitoring Subscriber via Go RPC)

-------------------------------------

## Running the Program

Step 1: Ensure VM2, VM3, and the Subscriber processes are active and reachable.

Step 2: Start the Web Server

go run main.go

-------------------------------------

## Example Session

[Monitor] Current Memory: 2 MB
[Monitor] Current Memory: 2 MB
[Monitor] Current Memory: 100 MB
[Monitor] Current Memory: 300 MB
[Monitor] Current Memory: 350 MB
!!! ALERT: High Memory Usage: 350 MB !!!

-------------------------------------

## Notes

- The continuous memory monitor ticks exactly every 5 seconds using `runtime.MemStats`.
- The memory alert threshold is hardcoded to a strict limit of 300 MB.
- Outbound event notifications are executed concurrently inside dedicated Goroutines to prevent blocking incoming user HTTP traffic.

-------------------------------------

## Packages Used

net
net/http
net/rpc
net/rpc/jsonrpc
runtime
time
encoding/base64
fmt
log