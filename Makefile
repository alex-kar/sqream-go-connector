generate_grpc:
	cd ./src/proto/ \
	protoc --go_out=stubs *.proto --go-grpc_out=stubs
run:
	go run ./cmd/sqream-connector/main.go
