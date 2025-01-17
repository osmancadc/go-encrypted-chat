#!/bin/bash

get_timestamp() {
    local timestamp=$(date +%Y-%m-%dT%H:%M:%S)
    local timezone=$(date +%z)
    local timezone_with_colon="${timezone:0:3}:${timezone:3:2}"
    echo "$timestamp$timezone_with_colon"
}

cd "$(dirname "$0")/.."

echo "$(get_timestamp) [INFO] Initializing the project go-encrypted-chat..."

echo "$(dirname "$0")"

echo "$(get_timestamp) [INFO] Executing unit tests..."
go test -coverprofile=coverage.out ./...

if [ $? -eq 0 ]; then
    echo "All the tests were successfull."

    echo "$(get_timestamp) [INFO] Building application..."
    go build -o go-encrypted-chat ./cmd/server

    echo "$(get_timestamp) [INFO] Starting application..."
    ./go-encrypted-chat
else
    echo "The test failed, please review the code."
    exit 1
fi
