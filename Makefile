# GOLANG with gRPC

## Envs
GOCMD=go


## Targets
.PHONY: run-grpc-server
run-grpc-server:
	cd cmd/grpc-server && $(GOCMD) run main.go

.PHONY: tidy
tidy:
	$(GOCMD) mod tidy


## Tests
.PHONY: test
test:
	$(GOCMD) test -v ./...

.PHONY: test-cover
cover:
	$(GOCMD) test -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

.PHONY: banchmark
banchmark:
	$(GOCMD) test -bench=. ./... -benchmem

.PHONY: evans
evans:
	evans -r repl --host localhost --port 50051

## gRPC
.PHONY: proto-all
proto-all:
	protoc -I ./proto --go_out=. --go-grpc_out=. ./proto/*.proto

.PHONY: grpcui
grpcui:
	grpcui -plaintext localhost:50051
