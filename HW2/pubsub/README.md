# Dedicated Monitoring Subscriber

This project implements a standalone monitoring node designed to asynchronously capture, format, and log real-time cluster health events dispatched across the network.

Processes:

1. Monitoring Subscriber Process

The Subscriber process binds to a telemetry port, accepts push alerts from the web cluster, and formats telemetry variables for human operator review.

-------------------------------------

## Supported Operations

NotificationService.ReceiveEvent   Accepts structured JSON memory alerts over RPC.

-------------------------------------

## Response Format / Log Layout

Success:

--- [New Event Received] ---
{
  "event_type": "HIGH_MEMORY_USAGE",
  "service": "web-server",
  "memory_mb": 350,
  "threshold_mb": 300,
  "timestamp": "2026-05-27T00:42:13+03:30"
}
----------------------------

-------------------------------------

## Network Ports Used

0.0.0.0:50053   (Listens across all active network interfaces for push metrics)

-------------------------------------

## Running the Program

Step 1: Boot up the Monitoring listener before triggering leaks on the Web Server.

go run subscriber.go

-------------------------------------

## Example Session

2026/05/27 00:40:29 Dedicated Monitoring Subscriber is listening on port 50053...

2026/05/27 00:42:13 
--- [New Event Received] ---
{
  "event_type": "HIGH_MEMORY_USAGE",
  "service": "web-server",
  "memory_mb": 350,
  "threshold_mb": 300,
  "timestamp": "2026-05-27T00:42:13+03:30"
}
----------------------------

-------------------------------------

## Notes

- Uses an elegant `json.MarshalIndent` mapping tool to render perfectly tabbed console printouts.
- Adheres cleanly to microservice design paradigms by insulating cluster processing nodes from rendering logs or formatting JSON strings locally.

-------------------------------------

## Packages Used

net
net/rpc
encoding/json
log