package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

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

type DataItem struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Version   int    `json:"version"`
	UpdatedBy string `json:"updated_by"`
}

func main() {
	fmt.Println("Starting automated Distributed Systems HW3 tests...")
	_ = os.MkdirAll("../results", 0755)

	runScenario1()
	runScenario2()
	runScenario3()
	runScenario4()

	fmt.Println("\nAll test pipelines finished! Check '../results/' for full metrics text logs.")
}

func runScenario1() {
	fmt.Println("-> Testing Scenario 1: Temporary Inconsistency")
	file, _ := os.Create("../results/scenario1.txt")
	defer file.Close()

	fmt.Fprintln(file, "=== SCENARIO 1: TEMPORARY INCONSISTENCY LOGS ===")
	key := "scen1_key"
	val := "100"

	reqBody, _ := json.Marshal(PutRequest{Key: key, Value: val, Consistency: "eventual", DelayMs: 1500})
	
	putStart := time.Now()
	resp, err := http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(reqBody))
	putLatency := time.Since(putStart)

	if err != nil {
		fmt.Fprintf(file, "Error running test: R1 offline: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var pResp PutResponse
	json.NewDecoder(resp.Body).Decode(&pResp)

	fmt.Fprintf(file, "PUT Latency: %v\n", putLatency)
	fmt.Fprintf(file, "Replicas Updated Immediately (Ack): %d\n", pResp.UpdatedReplicas)

	staleReads := 0
	getStart := time.Now()
	getResp, err := http.Get("http://localhost:8082/get?key=" + key)
	getLatency := time.Since(getStart)

	fmt.Fprintf(file, "Immediate GET Latency from Replica 2: %v\n", getLatency)
	if err != nil || getResp.StatusCode == http.StatusNotFound {
		staleReads++
		fmt.Fprintln(file, "Immediate GET status: Stale Read Detected (Data not present yet)")
	}
	if getResp != nil { getResp.Body.Close() }

	time.Sleep(2500 * time.Millisecond)

	getResp2, err := http.Get("http://localhost:8082/get?key=" + key)
	if err == nil && getResp2.StatusCode == http.StatusOK {
		var item DataItem
		json.NewDecoder(getResp2.Body).Decode(&item)
		fmt.Fprintf(file, "Post-convergence GET status: Value is now '%s' (Successfully Converged)\n", item.Value)
		getResp2.Body.Close()
	}
	fmt.Fprintf(file, "Total Stale Reads counted: %d\n", staleReads)
}

func runScenario2() {
	fmt.Println("-> Testing Scenario 2: Replica Failure Impact Analysis")
	file, _ := os.Create("../results/scenario2.txt")
	defer file.Close()

	fmt.Fprintln(file, "=== SCENARIO 2: REPLICA FAILURE EXPERIMENTS ===")
	fmt.Fprintln(file, "[Tip] For manual checking: kill Replica 3 (port 8083) process and rerun this script.")

	key := "scen2_key"
	val := "alive_test"

	reqBodyE, _ := json.Marshal(PutRequest{Key: key, Value: val, Consistency: "eventual", DelayMs: 0})
	startE := time.Now()
	respE, errE := http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(reqBodyE))
	
	if errE == nil {
		var pResp PutResponse
		json.NewDecoder(respE.Body).Decode(&pResp)
		fmt.Fprintf(file, "Eventual Consistency: Latency=%v, Local Ack Updated Replicas=%d\n", time.Since(startE), pResp.UpdatedReplicas)
		respE.Body.Close()
	}

	reqBodyS, _ := json.Marshal(PutRequest{Key: key, Value: val+"_strong", Consistency: "strong", DelayMs: 0})
	startS := time.Now()
	respS, errS := http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(reqBodyS))

	if errS == nil {
		var pResp PutResponse
		json.NewDecoder(respS.Body).Decode(&pResp)
		fmt.Fprintf(file, "Strong Consistency: Response Code=%d, Msg='%s', Latency=%v, Replicas Updated=%d\n", respS.StatusCode, pResp.Message, time.Since(startS), pResp.UpdatedReplicas)
		respS.Body.Close()
	} else {
		fmt.Fprintf(file, "Strong Consistency Call failed on network constraint: %v\n", errS)
	}
}

func runScenario3() {
	fmt.Println("-> Testing Scenario 3: Concurrent Conflict Resolution (LWW)")
	file, _ := os.Create("../results/scenario3.txt")
	defer file.Close()

	fmt.Fprintln(file, "=== SCENARIO 3: CONCURRENT WRITES CONFLICT RESOLUTION ===")
	key := "conflict_key"

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		req, _ := json.Marshal(PutRequest{Key: key, Value: "Value_From_Replica_1", Consistency: "eventual", DelayMs: 0})
		_, _ = http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(req))
	}()

	go func() {
		defer wg.Done()
		req, _ := json.Marshal(PutRequest{Key: key, Value: "Value_From_Replica_2", Consistency: "eventual", DelayMs: 0})
		_, _ = http.Post("http://localhost:8082/put", "application/json", bytes.NewBuffer(req))
	}()

	wg.Wait()
	time.Sleep(1500 * time.Millisecond)

	for _, port := range []string{"8081", "8082", "8083"} {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/get?key=%s", port, key))
		if err == nil && resp.StatusCode == http.StatusOK {
			var item DataItem
			json.NewDecoder(resp.Body).Decode(&item)
			fmt.Fprintf(file, "Replica Node [:%s] Converged State: Value=%s, Version=%d, UpdatedBy=%s\n", port, item.Value, item.Version, item.UpdatedBy)
			resp.Body.Close()
		}
	}
}

func runScenario4() {
	fmt.Println("-> Testing Scenario 4: Network Latency Impact Benchmarks")
	file, _ := os.Create("../results/scenario4.txt")
	defer file.Close()

	fmt.Fprintln(file, "=== SCENARIO 4: METRICS MATRIX UNDER ARTIFICIAL LATENCIES ===")

	delays := []int{0, 500, 2000}
	models := []string{"eventual", "strong"}

	for _, model := range models {
		for _, delay := range delays {
			key := fmt.Sprintf("scen4_key_%s_%d", model, delay)
			val := "metric_payload_data"

			reqBody, _ := json.Marshal(PutRequest{Key: key, Value: val, Consistency: model, DelayMs: delay})
			
			wallClockStart := time.Now()
			putStart := time.Now()
			resp, err := http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(reqBody))
			putLatency := time.Since(putStart)

			if err != nil {
				fmt.Fprintf(file, "Model: %s | Delay: %dms -> Connection Error on PUT: %v\n", model, delay, err)
				continue
			}
			resp.Body.Close()

			staleReads := 0
			var convergenceTime time.Duration
			isConverged := false
			
			// ثبت زمان شروع حلقه برای جلوگیری از گیر کردن ابدی (Timeout)
			loopStart := time.Now() 
			hasTimedOut := false

			for !isConverged {
				// اگر بیش از ۵ ثانیه گذشت و سیستم همگرا نشد (مثلاً به خاطر نود خاموش)
				if time.Since(loopStart) > 5*time.Second {
					hasTimedOut = true
					break
				}

				healthyPeers := 0
				for _, peerPort := range []string{"8082", "8083"} {
					gResp, gErr := http.Get(fmt.Sprintf("http://localhost:%s/get?key=%s", peerPort, key))
					if gErr == nil && gResp.StatusCode == http.StatusOK {
						var item DataItem
						json.NewDecoder(gResp.Body).Decode(&item)
						if item.Value == val {
							healthyPeers++
						} else {
							staleReads++
						}
						gResp.Body.Close()
					} else {
						staleReads++
					}
				}

				if healthyPeers == 2 {
					isConverged = true
					convergenceTime = time.Since(wallClockStart)
				} else {
					time.Sleep(40 * time.Millisecond)
				}
			}

			if hasTimedOut {
				fmt.Fprintf(file, "Model: %-8s | Delay Setting: %4dms | PUT Latency: %11v | Convergence Time: %11s | Stale Reads Checked: %2d\n",
					model, delay, putLatency, "TIMEOUT (Node Offline)", staleReads)
			} else {
				fmt.Fprintf(file, "Model: %-8s | Delay Setting: %4dms | PUT Latency: %11v | Convergence Time: %11v | Stale Reads Checked: %2d\n",
					model, delay, putLatency, convergenceTime, staleReads)
			}
		}
	}
}