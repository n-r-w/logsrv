.PHONY: build test run runbuild codegen rebuild tidy race docker-up docker-up-d docker-down docker-attach

build:	
	go build -v -o . ./cmd/logsrv.go

rebuild:	
	go build -a -v -o . ./cmd/logsrv.go

race:	
	go run -race ./cmd/logsrv.go -config-path ./config.toml

run:	
	go run ./cmd/logsrv.go -config-path ./config.toml

runbuild:
	./bin/logsrv

codegen:	
	protoc --go_out=./internal/presenter/grpc/generated --go_opt=paths=source_relative --go-grpc_out=./internal/presenter/grpc/generated --go-grpc_opt=paths=source_relative proto/logsrv.proto	
	wire ./internal/di	

tidy:
	go mod tidy

docker-up:
	docker-compose up 

docker-up-d:
	docker-compose up -d 

docker-down:
	docker-compose down

docker-attach:
	docker attach --sig-proxy=false logsrv-logsrv-1

.DEFAULT_GOAL := run
