# Authentication RPC Service

This project implements a decoupled authentication service that safely isolates sensitive user profile databases from public-facing server network tiers.

Processes:

1. Auth Service Process

The Auth Service receives raw credentials via high-performance JSON-RPC endpoints, parses an independent data store file, and securely verifies accounts.

-------------------------------------

## Supported Operations

AuthService.Login   Verifies Username and Password against local records.

-------------------------------------

## Response Format

Success:

{ "Success": true, "Message": "Authenticated" }

Errors:

{ "Success": false, "Message": "Invalid credentials" }

-------------------------------------

## Local Storage Files Used

users.json     Contains a key-value mapping of authorized system accounts.

Example schema:
{
  "admin": "admin123",
  "alice": "alice123",
  "bob": "bob123"
}

-------------------------------------

## Running the Program

Step 1: Populate users.json in the execution directory.

Step 2: Start the Auth Service

go run main.go

-------------------------------------

## Example Session

2026/05/27 00:40:24 Auth Service (VM2) JSON-RPC server is running on port 50051...
2026/05/27 00:41:27 [+] SUCCESS login: user='bob'
2026/05/27 00:43:02 [-] FAILED login: user='unknown_user'

-------------------------------------

## Notes

- Implements a modern `jsonrpc.ServeConn` handler rather than native Go binary formats to support language-agnostic integration patterns.
- Includes a robust error-handling fallback configuration that sets up basic temporary accounts if `users.json` is missing or corrupted.

-------------------------------------

## Packages Used

net
net/rpc
net/rpc/jsonrpc
encoding/json
io/ioutil
log