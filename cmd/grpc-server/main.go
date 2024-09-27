package main

import (
	"context"
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	userDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	pb "github.com/carlosdbarros/go-grpc-user-manage/internal/pb/user"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
)

func main() {
	db, err := configs.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create route handlers
	repo := database.NewUserDB(db)
	userHandler := NewUserHandler(repo)

	// create a gRPC server instance
	serverOpts := []grpc.ServerOption{}
	serverOpts = append(serverOpts, grpc.StreamInterceptor(StreamServerInterceptorCustom()))
	server := grpc.NewServer(serverOpts...)

	// register the service intances with the grpc server
	pb.RegisterUserServiceServer(server, userHandler)
	reflection.Register(server)

	// create a TCP listener on the specified port
	const addr = "0.0.0.0:50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// start listening to requests
	log.Printf("🚀 Server listening on %v", addr)
	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// StreamServerInterceptorCustom
func StreamServerInterceptorCustom() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Printf("StreamServerInterceptorCustom => %v", info.FullMethod)
		return handler(srv, stream)
	}
}

// UnaryServerInterceptorCustom
// func UnaryServerInterceptorCustom() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		log.Printf("UnaryServerInterceptorCustom => %v", req)
// 		return handler(ctx, req)
// 	}
// }

type UserHandler struct {
	pb.UnimplementedUserServiceServer

	Repo userDomain.UserRepository
}

func NewUserHandler(repo userDomain.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) CreateUser(_ context.Context, input *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := userDomain.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err = h.Repo.AddUser(user)
	if err != nil {
		log.Printf("Failed to add user: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	//log.Printf("Successfully created user: %v", user)
	return &pb.CreateUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *UserHandler) CreateUserAddress(_ context.Context, input *pb.CreateUserAddressRequest) (*pb.CreateUserAddressResponse, error) {
	//log.Printf("CreateUserAddress: %v", input)
	addressInput := make([]*userDomain.Address, 0)
	for _, a := range input.Addresses {
		address, err := userDomain.NewAddress(a.Street, a.Number, a.Complement, a.City, a.State, a.Country, a.ZipCode)
		if err != nil {
			//log.Printf("Failed to create address: %v", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		addressInput = append(addressInput, address)
	}
	userAddress, err := userDomain.NewUserAddress(input.Name, input.Emails, input.Phones, addressInput)
	if err != nil {
		//log.Printf("Failed to create user address: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	addressResponse := make([]*pb.Address, 0, 10)
	for _, a := range userAddress.Addresses {
		addressResponse = append(addressResponse, &pb.Address{
			Street:     a.Street,
			Number:     a.Number,
			Complement: a.Complement,
			City:       a.City,
			State:      a.State,
			Country:    a.Country,
			ZipCode:    a.ZipCode,
		})
	}
	return &pb.CreateUserAddressResponse{
		Name:      userAddress.Name,
		Emails:    userAddress.Emails,
		Phones:    userAddress.Phones,
		Addresses: addressResponse,
	}, nil
}

func (h *UserHandler) CreateUserStreamStream(stream pb.UserService_CreateUserStreamStreamServer) error {
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
		user, err = h.Repo.AddUser(user)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		err = stream.Send(&pb.CreateUserResponse{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			//log.Printf("Failed to send user: %v", err)
			return err
		}
		//log.Printf("Successfully created user: %v", user)
	}
}
