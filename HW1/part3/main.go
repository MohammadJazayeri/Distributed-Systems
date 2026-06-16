package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type ComputeResponse struct {
	Operation string   `json:"operation,omitempty"`
	A         *float64 `json:"a,omitempty"`
	B         *float64 `json:"b,omitempty"`
	Result    *float64 `json:"result,omitempty"`
	Error     string   `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method_not_allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func computeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method_not_allowed",
		})
		return
	}

	q := r.URL.Query()
	op := strings.TrimSpace(q.Get("op"))
	aStr := strings.TrimSpace(q.Get("a"))
	bStr := strings.TrimSpace(q.Get("b"))

	if op == "" || aStr == "" || bStr == "" {
		writeJSON(w, http.StatusBadRequest, ComputeResponse{
			Error: "missing_parameters",
		})
		return
	}

	a, errA := strconv.ParseFloat(aStr, 64)
	b, errB := strconv.ParseFloat(bStr, 64)
	if errA != nil || errB != nil || math.IsNaN(a) || math.IsNaN(b) || math.IsInf(a, 0) || math.IsInf(b, 0) {
		writeJSON(w, http.StatusBadRequest, ComputeResponse{
			Error: "non_numeric_input",
		})
		return
	}

	var result float64
	switch op {
	case "ADD":
		result = a + b
	case "SUB":
		result = a - b
	case "MUL":
		result = a * b
	case "DIV":
		if b == 0 {
			writeJSON(w, http.StatusBadRequest, ComputeResponse{
				Operation: op,
				A:         &a,
				B:         &b,
				Error:     "division_by_zero",
			})
			return
		}
		result = a / b
	case "MAX":
		if a > b {
			result = a
		} else {
			result = b
		}
	default:
		writeJSON(w, http.StatusBadRequest, ComputeResponse{
			Error: "invalid_op",
		})
		return
	}

	writeJSON(w, http.StatusOK, ComputeResponse{
		Operation: op,
		A:         &a,
		B:         &b,
		Result:    &result,
	})
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/compute", computeHandler)

	port := "8080"
	log.Printf("server listening on :%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
