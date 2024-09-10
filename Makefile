# GOLANG with gRPC

## Variables
GOCMD=go

## Targets
.PHONY: run-grpc-server
run-grpc-server:
	$(GOCMD) run cmd/grpc-server/main.go

.PHONY: run-web-server
run-web-server:
	$(GOCMD) run cmd/server/main.go

.PHONY: tidy
tidy:
	$(GOCMD) mod tidy

.PHONY: test
test:
	$(GOCMD) test -v ./...

.PHONY: banchmark
banchmark:
	$(GOCMD) test -bench=. ./... -benchmem

## gRPC
.PHONY: proto-user
proto-user:
	protoc --go_out=. --go-grpc_out=. proto/user.proto

.PHONY: proto-permission
proto-permission:
	protoc --go_out=. --go-grpc_out=. proto/permission.proto

.PHONY: proto-all
proto-all:
	protoc -I ./proto --go_out=. --go-grpc_out=. ./proto/*.proto

.PHONY: grpcui
grpcui:
	grpcui -plaintext localhost:50051