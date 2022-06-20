.PHONY: build test run runbuild codegen rebuild tidy race docker-up docker-up-d docker-down docker-attach

build:	
	go build -v -o . ./cmd/logsrv.go

rebuild:
	wire ./internal/di
	go build -a -v -o . ./cmd/logsrv.go

race:
	wire ./internal/di
	go run -race ./cmd/logsrv.go -config-path ./config.toml

run:	
	go run ./cmd/logsrv.go -config-path ./config.toml

runbuild:
	./bin/logsrv

codegen:
	wire ./internal/di

tidy:
	go mod tidy

docker-up:
	docker-compose up --build

docker-up-d:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-attach:
	docker attach --sig-proxy=false logsrv-logsrv-1

.DEFAULT_GOAL := run
