package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "gogrpc/proto/user"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	usersList map[string]*pb.User
}

// GetUser — bidirectional streaming
// Client sends user IDs, server streams back user details
func (uS *UserService) GetUser(stream pb.UserService_GetUserServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		log.Printf("Got request for user details with ID: %s", req.Id)

		user, exist := uS.usersList[req.Id]
		if !exist {
			log.Printf("No user found for ID: %s", req.Id)
			continue
		}

		if err := stream.Send(user); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// CreateUser — bidirectional streaming
// Client sends new users, server streams back the created user with an ID
func (uS *UserService) CreateUser(stream pb.UserService_CreateUserServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		log.Printf("Received request to create user: %s (%s)", req.Name, req.Email)

		userID := uuid.NewString()
		req.Id = userID
		uS.usersList[userID] = req

		// Send back the created user immediately
		if err := stream.Send(req); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userService := &UserService{
		usersList: make(map[string]*pb.User),
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	fmt.Printf("✅ gRPC server started on port %s...\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
