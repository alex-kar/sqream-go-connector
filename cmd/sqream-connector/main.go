package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/alex-kar/sqream-go-connector/src/proto/stubs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

const (
	authAddr string = "localhost:9091"
	qhAddr   string = "localhost:9092"
)

type ResultChunkHeader struct {
	ColSzs []string
	Rows   int32
}

func main() {
	// Establish connection
	conn, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	log.Printf("Connected to the server")
	defer conn.Close()

	// Get Token
	authClient := pb.NewAuthenticationServiceClient(conn)
	request := pb.AuthRequest{
		AuthType:    pb.AuthenticationType_AUTHENTICATION_TYPE_IDP,
		User:        "user",
		Password:    "pass",
		AccessToken: "Z2hZaDdyMmhEWHFkdGJBN3c4em9SSndjcVBXQjI5a05XZjRHSHU4X1B0R1RmbzYzYm53NENGaUVIMGlSX0lLRjVhMUQ3c3JfbHQyVGtfRk1md3U5T1M3aXNlZlcwS2l4",
	}
	response, err := authClient.Auth(context.Background(), &request)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	if response.GetError() != nil {
		log.Fatalf("Authentication response has error: %v", response.GetError())
	}
	token := response.GetToken()

	// Add Bearer token to the header
	md := metadata.Pairs("authorization", "Bearer "+token)
	context := metadata.NewOutgoingContext(context.Background(), md)
	fmt.Println("Exchanged AccessToken with JWT")

	// Open Session
	openSessionRequest := pb.SessionRequest{
		TenantId: "tenant",
		Database: "master",
		SourceIp: "127.0.0.1",
	}
	session, err := authClient.Session(context, &openSessionRequest)
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}
	if session.GetError() != nil {
		log.Fatalf("Session response has error: %v", session.GetError())
	}
	contextId := session.GetContextId()

	// Establish second connection to QueryHanler TODO: Remove it once configured secured connection properly through Ambassador
	conn.Close()
	maxMsgSize := 1024 * 1024 * 256 // 256 MBytes
	conn, err = grpc.Dial(qhAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	log.Printf("Connected to the server")

	qhClient := pb.NewQueryHandlerServiceClient(conn)

	sql := "select * from t;"

	// Compile
	compilationRequest := pb.CompileRequest{
		ContextId:    contextId,
		Sql:          []byte(sql),
		Encoding:     "UTF-8",
		QueryTimeout: 0,
	}
	compilationResult, err := qhClient.Compile(context, &compilationRequest)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}
	if compilationResult.GetError() != nil {
		log.Fatalf("Compilation result has error: %v", compilationResult.GetError())
	}
	columnMetadata := compilationResult.GetColumns()
	for col, index := range columnMetadata {
		log.Printf("Column %v: %v", index, col)
	}
	contextId = compilationResult.GetContextId()
	queryType := compilationResult.GetQueryType()
	fmt.Println("Query type", queryType)

	// Execute
	executeRequest := pb.ExecuteRequest{
		ContextId: contextId,
	}
	execResponse, err := qhClient.Execute(context, &executeRequest)
	if err != nil {
		log.Fatalf("Failed to launch query: %v", err)
	}
	if execResponse.GetError() != nil {
		log.Fatalf("Execution response has error: %v", execResponse.GetError())
	}

	// Get status
	statusRequest := pb.StatusRequest{
		ContextId: contextId,
	}
	for {
		statusResponse, err := qhClient.Status(context, &statusRequest)
		if err != nil || statusResponse.GetError() != nil {
			fmt.Println("Query exeuction failed. Status=", statusResponse.GetError())
			break
		}
		if statusResponse.GetStatus() != 1 && statusResponse.GetStatus() != 6 { // 'running' or 'queued'
			fmt.Println("Query exeuction completed. Status=", statusResponse.GetStatus())
			break
		}
	}

	// Fetch result if "QueryType_QUERY_TYPE_QUERY" type
	fmt.Println("About to fetch result")
	fetchRequest := pb.FetchRequest{
		ContextId: contextId,
	}
	if queryType == 1 {
		for {
			fetchResponse, err := qhClient.Fetch(context, &fetchRequest)
			if err != nil {
				log.Fatalf("Failed to fetch results: %v", err)
			}
			if fetchResponse.GetError() != nil {
				log.Fatalf("Fetch result has error: %v", execResponse.GetError())
			}
			if fetchResponse.GetRetryFetch() == true {
				fmt.Println("Fetch retry")
				continue
			}
			header, err := parseHeader(fetchResponse.GetQueryResult())
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

	// Close statement
	closeRequest := pb.CloseRequest{
		ContextId: contextId,
	}
	closeStmtReq := pb.CloseStatementRequest{
		CloseRequest: &closeRequest,
	}
	closeStmtResponse, err := qhClient.CloseStatement(context, &closeStmtReq)
	if err != nil {
		log.Fatalf("Close statement failed: %v", err)
	}
	if closeStmtResponse.CloseResponse.GetError() != nil {
		log.Fatalf("Close Statement has error: %v", compilationResult.GetError())
	}
	fmt.Println("Closed statement")

	// Close session
	closeSessionReq := pb.CloseSessionRequest{
		CloseRequest: &closeRequest,
	}
	closeSessionResponse, err := qhClient.CloseSession(context, &closeSessionReq)
	_ = closeSessionResponse
	if err != nil {
		log.Fatalf("Close session failed: %v", err)
	}
	if closeStmtResponse.CloseResponse.GetError() != nil {
		log.Fatalf("Close session has error: %v", compilationResult.GetError())
	}
	fmt.Println("Closed session")
}

func parseHeader(bytes []byte) (ResultChunkHeader, error) {
	if bytes == nil {
		return ResultChunkHeader{}, errors.New("Byte array is nil")
	}
	if len(bytes) < 8 {
		return ResultChunkHeader{}, errors.New("Byte arrays has less then 8 bytes")
	}
	headerLength := binary.LittleEndian.Uint64(bytes[:8])
	var result ResultChunkHeader
	json.Unmarshal(bytes[8:8+headerLength], &result)
	return result, nil
}
