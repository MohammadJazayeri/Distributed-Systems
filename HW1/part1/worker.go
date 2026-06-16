package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func isValidCommand(cmd string) bool {
	switch cmd {
	case "ADD", "SUB", "MUL", "DIV", "MAX":
		return true
	}
	return false
}

func main() {

	reqFIFO := "/tmp/pipe_req"
	respFIFO := "/tmp/pipe_resp"

	// create FIFOs if missing
	for _, fifo := range []string{reqFIFO, respFIFO} {
		if _, err := os.Stat(fifo); os.IsNotExist(err) {
			err := syscall.Mkfifo(fifo, 0666)
			if err != nil {
				log.Fatalf("Failed creating fifo %s: %v", fifo, err)
			}
		}
	}

	reqFile, err := os.OpenFile(reqFIFO, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Cannot open request pipe: %v", err)
	}
	defer reqFile.Close()

	respFile, err := os.OpenFile(respFIFO, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Cannot open response pipe: %v", err)
	}
	defer respFile.Close()

	scanner := bufio.NewScanner(reqFile)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)

		if len(parts) == 0 {
			fmt.Fprintln(respFile, "ERR empty_request")
			continue
		}

		cmd := parts[0]

		// 1️⃣ check command first
		if !isValidCommand(cmd) {
			fmt.Fprintln(respFile, "ERR unknown_command")
			continue
		}

		// 2️⃣ check argument count
		if len(parts) != 3 {
			fmt.Fprintln(respFile, "ERR wrong_argument_count")
			continue
		}

		// 3️⃣ parse numbers
		a, err1 := strconv.ParseFloat(parts[1], 64)
		b, err2 := strconv.ParseFloat(parts[2], 64)

		if err1 != nil || err2 != nil {
			fmt.Fprintln(respFile, "ERR non_numeric_input")
			continue
		}

		switch cmd {

		case "ADD":
			fmt.Fprintf(respFile, "OK %g\n", a+b)

		case "SUB":
			fmt.Fprintf(respFile, "OK %g\n", a-b)

		case "MUL":
			fmt.Fprintf(respFile, "OK %g\n", a*b)

		case "DIV":
			if b == 0 {
				fmt.Fprintln(respFile, "ERR division_by_zero")
			} else {
				fmt.Fprintf(respFile, "OK %g\n", a/b)
			}
		case "MAX":
			if a >= b {
				fmt.Fprintf(respFile, "OK %g\n", a)
			} else {
				fmt.Fprintf(respFile, "OK %g\n", b)
			}

		}

	}

	if err := scanner.Err(); err != nil {
		log.Println("Pipe read error:", err)
	}
}
