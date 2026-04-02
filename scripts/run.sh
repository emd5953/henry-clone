#!/bin/bash
# Load .env and start the server
set -a
source .env
set +a
exec go run cmd/server/main.go
