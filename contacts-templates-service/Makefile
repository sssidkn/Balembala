PROTO_FILE=proto/report.proto
OUT_DIR=pkg/api/report

SERVER_BIN=server
SERVER_MAIN=cmd/main.go

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

generate:
	protoc -I ./proto -I ./proto/google \
        --go_out=./${OUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=./${OUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=./${OUT_DIR} \
        --grpc-gateway_opt=paths=source_relative \
        ./${PROTO_FILE}

build:
	go build -o $(SERVER_BIN) $(SERVER_MAIN)

run:
	./$(SERVER_BIN)

migrate_up:
	 goose -dir migrations postgres "postgres://root:1234@localhost:5435/reportDB?sslmode=disable" up

migrate_down :
	 goose -dir migrations postgres "postgres://root:1234@localhost:5435/reportDB?sslmode=disable" down

