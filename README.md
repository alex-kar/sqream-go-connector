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
protoc --go_out=stubs *.proto
```
