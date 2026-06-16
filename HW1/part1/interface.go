package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main() {

	reqFIFO := "/tmp/pipe_req"
	respFIFO := "/tmp/pipe_resp"

	// Try opening request pipe in NONBLOCK mode
	reqFile, err := os.OpenFile(reqFIFO, os.O_WRONLY|syscall.O_NONBLOCK, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Worker is not running. Please start worker.go first.")
		return
	}
	defer reqFile.Close()

	// Now open response pipe normally
	respFile, err := os.OpenFile(respFIFO, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Cannot open response pipe:", err)
		return
	}
	defer respFile.Close()

	userReader := bufio.NewReader(os.Stdin)
	respReader := bufio.NewReader(respFile)

	for {

		fmt.Print("> ")

		input, err := userReader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		_, err = fmt.Fprintln(reqFile, input)
		if err != nil {
			fmt.Println("Worker connection lost.")
			return
		}

		resp, err := respReader.ReadString('\n')
		if err != nil {
			fmt.Println("Worker connection lost.")
			return
		}

		fmt.Print(resp)

	}
}
