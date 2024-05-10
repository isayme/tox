.PHONY: dev-local
dev-local:
	CONF_FILE_PATH=./config/example.config.yaml go run main.go local

@PHONY: dev-server
dev-server:
	CONF_FILE_PATH=./config/example.config.yaml go run main.go server

@PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/tunnel.proto
