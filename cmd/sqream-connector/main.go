package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"strings"
)

type Driver struct{}

func (driver Driver) Open(name string) (driver.Conn, error) {
	fmt.Printf("Driver.Open() with args [{%s}]\n", name)
	args := strings.Split(name, " ")
	addr := fmt.Sprintf("%s:%s", args[0], args[1])
	log.Printf("Trying to connect driver to [%s]", addr)
	client := GRPCClient{}
	client.Init(parseConnParams())
	err := client.OpenSession()
	if err != nil {
		return nil, err
	}
	return Conn{
		client: &client,
	}, nil
}

type ConnParam struct {
}

type Conn struct {
	client Client
}

func (conn Conn) Prepare(query string) (driver.Stmt, error) {
	log.Printf("Prepare statement [%s]", query)
	return Stmt{sql: query}, nil
}

func (conn Conn) Close() error {
	log.Print("Close conn")
	return nil
}

func (conn Conn) Begin() (driver.Tx, error) {
	return nil, nil
}

type Stmt struct {
	sql    string
	Client *grpc.ClientConn
}

func (stmt Stmt) Close() error {
	log.Print("Close stmt")
	return nil
}

func (stmt Stmt) NumInput() int {
	log.Print("NumInput")
	return 0
}

func (stmt Stmt) Exec(args []driver.Value) (driver.Result, error) {
	log.Print("Exec stmt")
	return nil, nil
}

func (stmt Stmt) Query(args []driver.Value) (driver.Rows, error) {
	log.Print("Query stmt")
	rows := Rows{Client: stmt.Client}
	return rows, nil
}

type Rows struct {
	Client *grpc.ClientConn
}

func (rows Rows) Columns() []string {
	log.Printf("Return result metadata")
	return []string{}
}

func (rows Rows) Close() error {
	return nil
}

func (rows Rows) Next(dest []driver.Value) error {
	return io.EOF
}

func init() {
	fmt.Println("Register driver")
	sql.Register("sqream", &Driver{})
}

func main() {
	host := "4-52.isqream.com"
	port := "443"
	token := "Z2hZaDdyMmhEWHFkdGJBN3c4em9SSndjcVBXQjI5a05XZjRHSHU4X1B0R1RmbzYzYm53NENGaUVIMGlSX0lLRjVhMUQ3c3JfbHQyVGtfRk1md3U5T1M3aXNlZlcwS2l4"

	sqInfo := fmt.Sprintf("host=%s port=%s token=%s", host, port, token)
	db, err := sql.Open("sqream", sqInfo)
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	log.Printf("Connected to the server")
	defer db.Close()

	fmt.Println("About to ping server")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping server %v", err)
	}

	fmt.Println("About to connect to the server")
	conn, err := db.Conn(context.Background())
	if err != err {
		log.Fatalf("Failed to connect to the server: %v", err)
	}
	defer conn.Close()

	rows, err := db.Query("select 1;")
	if err != nil {
		log.Fatalf("Failed to execute statement: %v", err)
		return
	}

	for rows.NextResultSet() {
		//NOP
	}
}

func parseConnParams() *ConnParams {
	// return hardcoded values for now
	return &ConnParams{
		Host:  "4-52.isqream.com",
		Port:  443,
		Token: "Z2hZaDdyMmhEWHFkdGJBN3c4em9SSndjcVBXQjI5a05XZjRHSHU4X1B0R1RmbzYzYm53NENGaUVIMGlSX0lLRjVhMUQ3c3JfbHQyVGtfRk1md3U5T1M3aXNlZlcwS2l4",
	}
}
