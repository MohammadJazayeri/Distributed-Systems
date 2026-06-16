# Go IPC Calculator using Named Pipes

This project implements communication between two independent processes
using Linux named pipes (FIFO).

Processes:

1. Worker Process
2. Interface Process

The Interface reads commands from the user and sends them to the Worker.
The Worker processes the request and sends the result back.

-------------------------------------

## Supported Operations

ADD A B
SUB A B
MUL A B
DIV A B

A and B are floating-point numbers.

Example:

ADD 5.5 2
SUB 10 3
MUL 2 4.5
DIV 8 2

-------------------------------------

## Response Format

Success:

OK result

Example:

OK 7.5

Errors:

ERR unknown_command
ERR wrong_argument_count
ERR non_numeric_input
ERR division_by_zero
ERR empty_request

-------------------------------------

## Named Pipes Used

/tmp/pipe_req     (Interface -> Worker)
/tmp/pipe_resp    (Worker -> Interface)

-------------------------------------

## Running the Program

Step 1: Start the Worker

go run worker.go

Step 2: Start the Interface in another terminal

go run interface.go

-------------------------------------

## Example Session

> ADD 5 3
OK 8

> MUL 2.5 4
OK 10

> DIV 5 0
ERR division_by_zero

> TEST 1 2
ERR unknown_command

> ADD 1
ERR wrong_argument_count

-------------------------------------

## Notes

- Worker must be started before Interface.
- If Interface is started first, it will print:

Worker is not running. Please start worker.go first.

- Worker runs continuously and can handle multiple requests.
- Each request is processed independently.

-------------------------------------

## Packages Used

os
bufio
fmt
strings
strconv
syscall
log