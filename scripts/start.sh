#!/bin/bash

cd "$(dirname "$0")/.."


echo "Initializing the project go-encrypted-chat..."

echo "$(dirname "$0")"


echo "Executing unit tests..."
go test ./...

if [ $? -eq 0 ]; then
    echo "All the tests were successfull."

    echo "Building application..."
    go build -o go-encrypted-chat ./cmd/server

    echo "Starting application..."
    ./go-encrypted-chat
else
    echo "The test failed, please review the code."
    exit 1
fi
