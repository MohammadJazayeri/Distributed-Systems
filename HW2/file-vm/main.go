package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
)

type FileRequest struct{ FileName string }
type FileResponse struct{ Content []byte }
type FileService struct{}

type Event struct {
	Topic   string
	Message string
}

type NotificationService struct{}

func (n *NotificationService) ReceiveEvent(event Event, reply *bool) error {
	log.Printf("[PUB/SUB SUB] Received on Topic [%s]: %s", event.Topic, event.Message)
	*reply = true
	return nil
}

func (s *FileService) GetFile(req FileRequest, res *FileResponse) error {
	log.Printf("RPC File Request: %s", req.FileName)
	filePath := "./storage/" + req.FileName

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("File error: %v", err)
		return err
	}

	res.Content = data
	return nil
}

func main() {
	os.MkdirAll("./storage", 0755)

	fileSvc := new(FileService)
	rpc.Register(fileSvc)

	notifSvc := new(NotificationService)
	rpc.Register(notifSvc)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	log.Println("File RPC Service (VM3) is running on port 50052...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}