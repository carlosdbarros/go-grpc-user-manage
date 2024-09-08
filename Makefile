# GOLANG with gRPC

## Variables
GOCMD=go

## Targets
run:
	$(GOCMD) run cmd/grpc-server/main.go

tidy:
	$(GOCMD) mod tidy

.PHONY: test
test:
	$(GOCMD) test -v ./...

banchmark:
	$(GOCMD) test -bench=. ./... -benchmem

## gRPC
proto:
	protoc --go_out=. --go-grpc_out=. internal/proto/*.proto
