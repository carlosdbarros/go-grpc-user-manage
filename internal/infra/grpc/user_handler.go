package grpc

import (
	"context"
	userDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
	pb "github.com/carlosdbarros/go-grpc-user-manage/internal/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer

	Repo userDomain.UserRepository
}

func NewUserHandler(repo userDomain.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) CreateUser(_ context.Context, input *pb.CreateUserRequest) (*pb.User, error) {
	user, err := userDomain.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//user, err = h.Repo.AddUser(user)
	//if err != nil {
	//	log.Printf("Failed to add user: %v", err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *UserHandler) CreateUserStream(stream pb.UserService_CreateUserStreamStreamServer) error {
	for {
		input, err := stream.Recv()
		if err == io.EOF {
			//log.Printf("End of stream")
			return nil
		}
		if err != nil {
			//log.Printf("Failed to receive user: %v\nerro: %v", input, err)
			return err
		}
		user, err := userDomain.NewUser(input.Name, input.Email, input.Password)
		if err != nil {
			//log.Printf("Failed to create user: %v", err)
			return status.Error(codes.InvalidArgument, err.Error())
		}
		//user, err = h.Repo.AddUser(user)
		//if err != nil {
		//	return status.Error(codes.Internal, err.Error())
		//}
		err = stream.Send(&pb.CreateUserResponse{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			//log.Printf("Failed to send user: %v", err)
			return err
		}
		//log.Printf("User sent successfully => %v", user)
	}
}

func (h *UserHandler) FindAllUsers(_ context.Context, _ *pb.Empty) (*pb.FindAllUsersResponse, error) {
	users, err := h.Repo.FindAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &pb.FindAllUsersResponse{Users: pbUsers}, nil
}
