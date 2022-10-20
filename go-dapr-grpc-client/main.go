package main

import (
	"context"
	"encoding/json"
	pb "go-dapr-grpc-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

const (
	address = "localhost:3030"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println("Connected to server!")
	defer conn.Close()
	c := pb.NewTodoListClient(conn)

	for {
		time.Sleep(1 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", "server")
		r, err := c.GetTodolist(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("Received todo list with size: %v", r.GetSize())
		d, _ := json.Marshal(r.GetTodoLists())
		log.Printf("Data: %v", string(d))
	}
}
