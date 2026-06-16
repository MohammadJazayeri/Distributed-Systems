package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"runtime"
	"time"
)

type LoginRequest struct{ Username, Password string }
type LoginResponse struct{ Success bool; Message string }
type FileRequest struct{ FileName string }
type FileResponse struct{ Content []byte }

type Event struct {
	EventType   string `json:"event_type"`
	Service     string `json:"service"`
	MemoryMB    uint64 `json:"memory_mb"`
	ThresholdMB uint64 `json:"threshold_mb"`
	Timestamp   string `json:"timestamp"`
}

const (
	VM2_ADDR  = "192.168.56.104:50051"
	VM3_ADDR  = "192.168.56.102:50052"
	THRESHOLD = 300 * 1024 * 1024
)


func monitorMemory() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memMB := m.Alloc / 1024 / 1024
		
		fmt.Printf("[Monitor] Current Memory: %d MB\n", memMB)
		
		thresholdMB := uint64(THRESHOLD / 1024 / 1024)
		
		if memMB >= thresholdMB {
			fmt.Printf("!!! ALERT: High Memory Usage: %d MB !!!\n", memMB)
			
			event := Event{
				EventType:   "HIGH_MEMORY_USAGE",
				Service:     "web-server",
				MemoryMB:    memMB,
				ThresholdMB: thresholdMB,
				Timestamp:   time.Now().Format(time.RFC3339),
			}
			
			publish(event)
		}
		time.Sleep(5 * time.Second)
	}
}

var memoryHog [][]byte
func consumeMemoryHandler(w http.ResponseWriter, r *http.Request) {
	chunk := make([]byte, 50*1024*1024)
	for i := range chunk {
		chunk[i] = 1
	}
	memoryHog = append(memoryHog, chunk)

	fmt.Fprintf(w, "Consumed 50MB. Simulated memory leak successfully.")
}

func publish(event Event) {
	subscribers := []string{"192.168.56.104:50053"}

	for _, addr := range subscribers {
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			continue
		}

		go func(c *rpc.Client, ev Event) {
			defer c.Close()
			var rep bool
			c.Call("NotificationService.ReceiveEvent", ev, &rep)
		}(client, event)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `<html><body style="text-align:center; font-family:sans-serif;">
            <h2>Login to Distributed System</h2>
            <form action="/login" method="post" style="border:1px solid #ccc; display:inline-block; padding:20px;">
                Username: <input name="username" required><br><br>
                Password: <input type="password" name="password" required><br><br>
                <input type="submit" value="Login">
            </form></body></html>`)
		return
	}

	// VM1 -> VM2 via JSON-RPC
	conn2, err := net.Dial("tcp", VM2_ADDR)
	if err != nil {
		http.Error(w, "Critical Error: Auth Service (VM2) is offline.", 500)
		return
	}
	defer conn2.Close()

	client2 := jsonrpc.NewClient(conn2)

	req := LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	var res LoginResponse

	err = client2.Call("AuthService.Login", req, &res)
	if err != nil || !res.Success {
		w.Header().Set("Content-Type", "text/html")
		reason := res.Message
		if err != nil {
			reason = err.Error()
		}
		fmt.Fprintf(w, `<html><body style="text-align:center; color:red;">
            <h2>Login Failed!</h2>
            <p>Reason: %s</p>
            <a href="/login">Return to Login Page</a>
        </body></html>`, reason)
		return
	}

	// VM1 -> VM3 via normal Go RPC
	client3, err := rpc.Dial("tcp", VM3_ADDR)
	if err != nil {
		fmt.Fprintf(w, "<h1>Login Success!</h1><p>Welcome %s. (But VM3 is down for image fetch)</p>", req.Username)
		return
	}
	defer client3.Close()

	fRes := FileResponse{}
	err = client3.Call("FileService.GetFile", FileRequest{"image.png"}, &fRes)

	var imgBase64Str string
	if err == nil {
		imgBase64Str = base64.StdEncoding.EncodeToString(fRes.Content)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body style="text-align:center;">
        <h1 style="color:green;">Access Granted!</h1>
        <p>Welcome, <b>%s</b></p>
        <hr>
        <h3>Protected Image from VM3:</h3>`+
		func() string {
			if imgBase64Str != "" {
				return fmt.Sprintf(`<img src="data:image/png;base64,%s" width="450" style="border:5px solid #555;" />`, imgBase64Str)
			}
			return `<p style="color:red;">Failed to load image from VM3 storage.</p>`
		}()+`
        <br><br>
        <a href="/login">Logout</a>
    </body></html>`, req.Username)
}

func main() {
	go monitorMemory()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consume-memory", consumeMemoryHandler)

	fmt.Println("Web Server running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
