#!/bin/bash
set -e

echo "Starting the Check for Backend..."

cd /backend
golangci-lint run
go mod download
go mod verify
go build ./...
go test ./...

echo "Done Successfully!"