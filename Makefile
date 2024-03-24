generate_grpc:
	cd ./src/proto/ \
	protoc --go_out=stubs *.proto --go-grpc_out=stubs
run:
	go run ./cmd/sqream-connector/main.go
run_debug:
	GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 go run ./cmd/sqream-connector/main.go
