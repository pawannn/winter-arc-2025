package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "gogrpc/proto/user"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userList = []pb.User{
	{
		Name:  "John Doe",
		Email: "test@gmail.com",
	},
	{
		Name:  "Jake",
		Email: "test2@gmail.com",
	},
	{
		Name:  "Anna",
		Email: "test3@gmail.com",
	},
}

func main() {
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Open a bi-directional stream
	stream, err := client.CreateUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			fmt.Printf("User created: ID=%s, Name=%s, Email=%s\n", resp.Id, resp.Name, resp.Email)
		}
		close(done)
	}()

	for i := range userList {
		userList[i].Id = uuid.NewString()
		if err := stream.Send(&userList[i]); err != nil {
			log.Fatalf("Error sending: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Error closing send: %v", err)
	}

	<-done
	fmt.Println("All users sent and responses received")
}
