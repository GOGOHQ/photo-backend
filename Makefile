SHELL := /bin/bash

.PHONY: run tidy docker-up docker-down

run:
	GO111MODULE=on go run ./cmd/server

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down 