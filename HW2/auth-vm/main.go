package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type LoginRequest struct {
	Username string
	Password string
}

type Event struct {
	Topic   string
	Message string
}

type LoginResponse struct {
	Success bool
	Message string
}

type AuthService struct{}
type NotificationService struct{}

func (n *NotificationService) ReceiveEvent(event Event, reply *bool) error {
	log.Printf("[PUB/SUB SUB] Received on Topic [%s]: %s", event.Topic, event.Message)
	*reply = true
	return nil
}

func loadUsers() map[string]string {
	content, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Println("Error reading users.json, using default admin.")
		return map[string]string{
			"admin": "admin123",
			"alice": "alice123",
			"bob":   "bob123",
		}
	}

	var users map[string]string
	if err := json.Unmarshal(content, &users); err != nil {
		log.Println("Invalid users.json, using default admin.")
		return map[string]string{
			"admin": "admin123",
			"alice": "alice123",
			"bob":   "bob123",
		}
	}
	return users
}

func (s *AuthService) Login(req LoginRequest, res *LoginResponse) error {
	users := loadUsers()
	savedPassword, exists := users[req.Username]

	if !exists || savedPassword != req.Password {
		log.Printf("[-] FAILED login: user='%s'", req.Username)
		res.Success = false
		res.Message = "Invalid credentials"
		return nil
	}

	log.Printf("[+] SUCCESS login: user='%s'", req.Username)
	res.Success = true
	res.Message = "Authenticated"
	return nil
}

func main() {
	authSvc := new(AuthService)
	if err := rpc.Register(authSvc); err != nil {
		log.Fatal("Register AuthService error:", err)
	}

	notifSvc := new(NotificationService)
	if err := rpc.Register(notifSvc); err != nil {
		log.Fatal("Register NotificationService error:", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	log.Println("Auth Service (VM2) JSON-RPC server is running on port 50051...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
