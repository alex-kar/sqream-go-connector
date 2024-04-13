package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	pb "github.com/alex-kar/sqream-go-connector/src/proto/stubs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

type Client interface {
	Init(params *ConnParams) error
	OpenSession() error
	Compile(sql string) error
	Execute() error
	Fetch() error
	CloseStmt() error
	CloseConn() error
}

type ConnParams struct {
	Host  string
	Port  int32
	Token string
}

type GRPCClient struct {
	ConnParams *ConnParams
	Conn       *grpc.ClientConn
	Token      string
	ContextId  string
	Columns    []*pb.ColumnMetadata
	QueryType  pb.QueryType
}

type Chunk struct {
}

func (client *GRPCClient) Init(params *ConnParams) error {
	log.Printf("Init gRPC client with connection params [params=%v]", params)
	client.ConnParams = params
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	maxMsgSize := 1024 * 1024 * 256 // 256 MBytes
	addr := fmt.Sprintf("%s:%d", params.Host, params.Port)
	log.Printf("Establish gRPC connection [%s]", addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return err
	}
	client.Conn = conn
	return nil
}

func (client *GRPCClient) OpenSession() error {
	// get JWT
	authClient := pb.NewAuthenticationServiceClient(client.Conn)
	authRequest := pb.AuthRequest{
		AuthType:    pb.AuthenticationType_AUTHENTICATION_TYPE_IDP,
		User:        "ignore",
		Password:    "ignore",
		AccessToken: client.ConnParams.Token,
	}
	log.Printf("Send auth request [%v]\n", &authRequest)
	response, err := authClient.Auth(context.Background(), &authRequest)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
		return err
	}
	if response.GetError() != nil {
		log.Fatalf("Authentication response has error: %v", response.GetError())
		return errors.New(response.GetError().String())
	}
	client.Token = response.GetToken()

	// open session
	clientInfo := pb.ClientInfo{
		Version:    "",
		SourceType: pb.SourceType_CLI,
	}
	openSessionReq := pb.SessionRequest{
		TenantId:   "tenant",
		Database:   "master",
		SourceIp:   "127.0.0.1",
		ClientInfo: &clientInfo,
		PoolName:   "sqream",
	}
	sessResp, err := authClient.Session(buildContext(client.Token), &openSessionReq)
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
		return err
	}
	if sessResp.GetError() != nil {
		log.Fatalf("Open session response has failure: %v", sessResp.GetError())
		return errors.New(sessResp.GetError().String())
	}
	client.ContextId = sessResp.ContextId
	return nil
}

func (client *GRPCClient) Compile(sql string) error {
	qhClient := pb.NewQueryHandlerServiceClient(client.Conn)
	openStmtReq := pb.CompileRequest{
		ContextId:    client.ContextId,
		Sql:          []byte(sql),
		Encoding:     "UTF-8",
		QueryTimeout: 0,
	}
	resp, err := qhClient.Compile(buildContext(client.Token), &openStmtReq)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
		return err
	}
	if resp.GetError() != nil {
		log.Fatalf("Compilation response has error: %v", resp.GetError())
		return errors.New(resp.GetError().String())
	}
	client.Columns = resp.GetColumns()
	client.QueryType = resp.GetQueryType()
	client.ContextId = resp.GetContextId()
	return nil
}

func (client *GRPCClient) Execute() error {
	qhClient := pb.NewQueryHandlerServiceClient(client.Conn)
	execReq := pb.ExecuteRequest{
		ContextId: client.ContextId,
	}
	// launch query
	resp, err := qhClient.Execute(buildContext(client.Token), &execReq)
	if err != nil {
		log.Fatalf("Qeury execution failed: %v", err)
		return err
	}
	if resp.GetError() != nil {
		log.Fatalf("Query execution response has error: %v", resp.GetError())
		return errors.New(resp.GetError().String())
	}
	client.ContextId = resp.GetContextId()
	// get status
	statusRequest := pb.StatusRequest{
		ContextId: client.ContextId,
	}
	for {
		statusResponse, err := qhClient.Status(buildContext(client.Token), &statusRequest)
		if err != nil {
			fmt.Println("Query exeuction failed. Status=", statusResponse.GetError())
			return err
		}
		if statusResponse.GetError() != nil {
			fmt.Println("Query exeuction failed. Status=", statusResponse.GetError())
			return errors.New(statusResponse.GetError().String())
		}
		if statusResponse.GetStatus() != 1 && statusResponse.GetStatus() != 6 { // 'running' or 'queued'
			fmt.Println("Query exeuction completed. Status=", statusResponse.GetStatus())
			break
		}
	}
	return nil
}

func (client *GRPCClient) Fetch() error {
	qhClient := pb.NewQueryHandlerServiceClient(client.Conn)
	fetchRequest := pb.FetchRequest{
		ContextId: client.ContextId,
	}
	if client.QueryType == 1 {
		for {
			fetchResponse, err := qhClient.Fetch(buildContext(client.Token), &fetchRequest)
			if err != nil {
				log.Fatalf("Failed to fetch results: %v", err)
			}
			if fetchResponse.GetError() != nil {
				log.Fatalf("Fetch result has error: %v", fetchResponse.GetError())
			}
			if fetchResponse.GetRetryFetch() == true {
				fmt.Println("Fetch retry")
				continue
			}
			chunkParser := new(ChunkParser)
			header, err := chunkParser.Header(fetchResponse.GetQueryResult())
			if err != nil {
				fmt.Println("Failed to parse chunk header", err)
				break
			}
			if header.Rows == 0 {
				fmt.Println("Loaded all chunks")
				break
			}
			fmt.Printf("Loaded chunk with [{%d}]\n", header.Rows)
		}
	}
	return nil
}

func (client *GRPCClient) CloseStmt() error {
	qhClient := pb.NewQueryHandlerServiceClient(client.Conn)
	closeRequest := pb.CloseRequest{
		ContextId: client.ContextId,
	}
	closeStmtReq := pb.CloseStatementRequest{
		CloseRequest: &closeRequest,
	}
	closeStmtResponse, err := qhClient.CloseStatement(buildContext(client.Token), &closeStmtReq)
	if err != nil {
		log.Fatalf("Close statement failed: %v", err)
	}
	if closeStmtResponse.CloseResponse.GetError() != nil {
		log.Fatalf("Close Statement has error: %v", closeStmtResponse.CloseResponse.GetError())
	}
	fmt.Println("Closed statement")
	return nil
}

func (client *GRPCClient) CloseConn() error {
	qhClient := pb.NewQueryHandlerServiceClient(client.Conn)
	closeRequest := pb.CloseRequest{
		ContextId: client.ContextId,
	}
	closeSessionReq := pb.CloseSessionRequest{
		CloseRequest: &closeRequest,
	}
	closeSessionResponse, err := qhClient.CloseSession(buildContext(client.Token), &closeSessionReq)
	_ = closeSessionResponse
	if err != nil {
		log.Fatalf("Close session failed: %v", err)
	}
	if closeSessionResponse.CloseResponse.GetError() != nil {
		log.Fatalf("Close session has error: %v", closeSessionResponse.CloseResponse.GetError())
	}
	fmt.Println("Closed session")
	return nil
}

func buildContext(token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(context.Background(), md)
}
