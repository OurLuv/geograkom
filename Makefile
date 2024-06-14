.PHONY: up

local: #local run
	go run cmd/main.go

up:
	docker-compose up -d --build

test:
	/usr/local/bin/go test -timeout 30s -run ^TestGetRouteById$ github.com/OurLuv/geograkom/internal/storage