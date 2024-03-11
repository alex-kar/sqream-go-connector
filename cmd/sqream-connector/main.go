package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/alex-kar/sqream-go-connector/src/proto/stubs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr string = "localhost:9090"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthenticationServiceClient(conn)
	request := pb.AuthRequest{
		AuthType:    pb.AuthenticationType_AUTHENTICATION_TYPE_IDP,
		User:        "user",
		Password:    "pass",
		AccessToken: "access_token",
	}
	response, err := c.Auth(context.Background(), &request)
	if err != nil {
		log.Fatalf("Failed to execute gRPC request: %v", err)
	}
	fmt.Printf("Response: %v\n", response)
}
