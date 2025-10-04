package main

import (
	"context"
	"fmt"
	pb "gogrpc/proto/user"
	"gogrpc/utils"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	usersList map[string]*pb.User
}

func (uS *UserService) GetUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	if user.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Please provide a valid id")
	}

	if _, exist := uS.usersList[user.Id]; !exist {
		return nil, status.Error(codes.NotFound, "No user found with the given userID")
	}

	userDetails := uS.usersList[user.Id]

	return userDetails, nil
}

func (uS *UserService) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	if user.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Please provide a valid id")
	}
	if !utils.ValidateName(user.Name) {
		return nil, status.Error(codes.InvalidArgument, "Please provide a valid name")
	}
	if !utils.ValidateEmail(user.Email) {
		return nil, status.Error(codes.InvalidArgument, "Please provide a valid email")
	}

	if _, exist := uS.usersList[user.Id]; exist {
		return nil, status.Error(codes.AlreadyExists, "The User ID already exist for a user")
	}

	uS.usersList[user.Id] = user

	return user, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserService{})

	fmt.Println("Started listening on TCP::8080...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
