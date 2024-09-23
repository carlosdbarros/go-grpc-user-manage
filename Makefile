# GOLANG with gRPC

## Envs
GOCMD=go
PB_PATH=./internal/pb


## Targets
.PHONY: run-grpc-server
run-grpc-server:
	cd cmd/grpc-server && \
	$(GOCMD) build -o grpc-server && \
	./grpc-server

.PHONY: run-http-server
run-http-server:
	cd cmd/http-server && \
	$(GOCMD) build -o http-server && \
	./http-server

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


.PHONY: evans
evans:
	evans -r repl --host localhost --port 50051

## gRPC
.PHONY: protogen
protogen:
	protoc -I ./proto --go_out=$(PB_PATH) --go_opt=paths=source_relative \
 		--go-grpc_out=$(PB_PATH) --go-grpc_opt=paths=source_relative \
 		--grpc-gateway_out=$(PB_PATH) --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
 		./proto/**/*.proto ./proto/**/**/*.proto

.PHONY: grpcui
grpcui:
	grpcui -plaintext localhost:50051
