#!/bin/bash

# Download dependencies
go mod download

# Build the application
go build -o app

# Run the application
./app 