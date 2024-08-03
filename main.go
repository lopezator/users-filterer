package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	pb "buf.build/gen/go/lopezator/filterer/grpc/go/lopezator/filterer/v1/filtererv1grpc"
	filtererpb "buf.build/gen/go/lopezator/filterer/protocolbuffers/go/lopezator/filterer/v1"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:1337", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close gRPC connection: %v", err)
		}
	}(conn)

	// Create a client for the FiltererService
	client := pb.NewFiltererServiceClient(conn)

	// Create a request payload
	req := &filtererpb.FilterRequest{
		Expr: "email == 'bob@example.com'",
	}

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make the gRPC call
	resp, err := client.Filter(ctx, req)
	if err != nil {
		log.Fatalf("Error making gRPC call: %v", err)
	}

	// Print the response
	fmt.Println("Response:", resp)

	// Set up a connection to the database
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/users?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	// Prepare the SQL query using the response from the gRPC call
	query := fmt.Sprintf("select * from users WHERE %s", resp.Where)
	for i := range resp.Args {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, "?", placeholder, 1)
	}

	// Convert args from []string to []interface{} for the query
	args := make([]interface{}, len(resp.Args))
	for i, v := range resp.Args {
		args[i] = v
	}

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Failed to close rows: %v", err)
		}
	}(rows)

	// Process the query results
	for rows.Next() {
		var id int
		var username, email string
		var createdAt time.Time
		if err := rows.Scan(&id, &username, &email, &createdAt); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("User: %d, %s, %s, %s\n", id, username, email, createdAt)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}
}
