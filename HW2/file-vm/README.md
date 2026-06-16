# Protected File Storage RPC Service

This project implements an isolation layer for static assets and sensitive media files within the secure virtual machine network space.

Processes:

1. File Service Process

The File Service serves internal system requests by retrieving data safely from an encrypted or protected directory and streaming it back to authorized clients.

-------------------------------------

## Supported Operations

FileService.GetFile   Reads a file from storage and returns its raw byte stream payload.

-------------------------------------

## Response Format

Success:

Returns a raw binary byte array (`[]byte`) mapping the exact data structure of the target asset.

Errors:

Returns a standard Go `error` object indicating file-read or directory-permission exceptions.

-------------------------------------

## Local Storage Directories Used

./storage/            (Root directory containing protected assets)
./storage/image.png   (Target dashboard asset pulled by Web Server)

-------------------------------------

## Running the Program

Step 1: Start the File Service to generate directory assets

go run main.go

Step 2: Place an image named `image.png` into the auto-generated `./storage/` folder.

-------------------------------------

## Example Session

2026/05/27 00:40:33 File RPC Service (VM3) is running on port 50052...
2026/05/27 00:41:27 RPC File Request: image.png

-------------------------------------

## Notes

- The service automatically instantiates the storage repository on its initial startup sequence with a `0755` file permissions mask.
- Avoids exposing HTTP bindings to end-users directly, completely cutting off public vectors for remote directory traversal attacks.

-------------------------------------

## Packages Used

net
net/rpc
os
io/ioutil
log