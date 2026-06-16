package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type DataItem struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Version   int    `json:"version"`
	UpdatedBy string `json:"updated_by"`
}

type Config struct {
	ID    string   `json:"id"`
	Port  string   `json:"port"`
	Peers []string `json:"peers"`
}

type PutRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Consistency string `json:"consistency"`
	DelayMs     int    `json:"delay_ms"`
}

type PutResponse struct {
	Message         string `json:"message"`
	UpdatedReplicas int    `json:"updated_replicas"`
	Version         int    `json:"version"`
}

var (
	config Config
	store  = make(map[string]DataItem)
	mu     sync.RWMutex
)

func main() {
	configFile := flag.String("config", "", "Path to the JSON config file")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Please specify a config file using -config flag")
	}

	file, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/replicate", handleReplicate)

	log.Printf("[%s] Server running on port %s...", config.ID, config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	mu.RLock()
	item, exists := store[key]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	var req PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	current, exists := store[req.Key]
	nextVersion := 1
	if exists {
		nextVersion = current.Version + 1
	}

	item := DataItem{
		Key:       req.Key,
		Value:     req.Value,
		Version:   nextVersion,
		UpdatedBy: config.ID,
	}
	store[req.Key] = item
	mu.Unlock()

	updatedCount := 1 

	if req.Consistency == "strong" {
		var wg sync.WaitGroup
		var lock sync.Mutex

		for _, peer := range config.Peers {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				if sendReplication(p, item, req.DelayMs) {
					lock.Lock()
					updatedCount++
					lock.Unlock()
				}
			}(peer)
		}
		wg.Wait()

		quorum := (len(config.Peers)+1)/2 + 1
		if updatedCount < quorum {
			mu.Lock()
			if exists {
				store[req.Key] = current
			} else {
				delete(store, req.Key)
			}
			mu.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(PutResponse{
				Message:         "Write failed: Quorum not reached",
				UpdatedReplicas: updatedCount,
				Version:         0,
			})
			return
		}
	} else {
		for _, peer := range config.Peers {
			go sendReplication(peer, item, req.DelayMs)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PutResponse{
		Message:         "Write successful",
		UpdatedReplicas: updatedCount,
		Version:         item.Version,
	})
}

func handleReplicate(w http.ResponseWriter, r *http.Request) {
	var incoming DataItem
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	current, exists := store[incoming.Key]
	if exists {
		if incoming.Version < current.Version {
			w.WriteHeader(http.StatusOK)
			return
		}
		if incoming.Version == current.Version {
			if incoming.UpdatedBy <= current.UpdatedBy {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	store[incoming.Key] = incoming
	w.WriteHeader(http.StatusOK)
}

func sendReplication(peerURL string, item DataItem, delayMs int) bool {
	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	payload, _ := json.Marshal(item)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(peerURL+"/replicate", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}