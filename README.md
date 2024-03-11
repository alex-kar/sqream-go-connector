# sqream-go-connector
SQream Golang Connector

## Generate gRPC code

Install plugin  
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Compile
```
cd src/proto
mkdir stubs
protoc --go_out=stubs *.proto --go-grpc_out=stubs
```

To install plugin
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
