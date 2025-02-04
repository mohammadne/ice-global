#! /bin/bash

# setup the test dependencies
docker compose -f ./deployments/docker/compose.local.yml up -d

# migrate the database
go run cmd/migrations/main.go --direction=up

# run the web-api
go run cmd/web-api/main.go --environment=local
