.PHONY: up

local: #local run
	go run cmd/main.go

up:
	docker-compose up -d --build
