package main

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
)

type Event struct {
	EventType   string `json:"event_type"`
	Service     string `json:"service"`
	MemoryMB    uint64 `json:"memory_mb"`
	ThresholdMB uint64 `json:"threshold_mb"`
	Timestamp   string `json:"timestamp"`
}

type NotificationService struct{}

func (n *NotificationService) ReceiveEvent(event Event, reply *bool) error {
	eventJSON, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return nil
	}

	log.Printf("\n--- [New Event Received] ---\n%s\n----------------------------", string(eventJSON))
	
	*reply = true
	return nil
}

func main() {
	notifSvc := new(NotificationService)
	rpc.Register(notifSvc)

	// گوش دادن روی تمام اینترفیس‌ها پورت 50053
	listener, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatal("Subscriber listen error:", err)
	}

	log.Println("Dedicated Monitoring Subscriber is listening on port 50053...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
