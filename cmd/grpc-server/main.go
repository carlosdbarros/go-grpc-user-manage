package main

import (
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	grpcInfra "github.com/carlosdbarros/go-grpc-user-manage/internal/infra/grpc"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/pb/user"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	db, err := configs.InitSqliteInMemory()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create route handlers
	repo := database.NewUserDB(db)
	userHandler := grpcInfra.NewUserHandler(repo)

	// create a gRPC server instance
	serverOpts := []grpc.ServerOption{}
	// serverOpts = append(serverOpts, grpc.UnaryInterceptor(UnaryServerInterceptorCustom()))
	// serverOpts = append(serverOpts, grpc.StreamInterceptor(StreamServerInterceptorCustom()))
	server := grpc.NewServer(serverOpts...)

	// register the service intances with the grpc server
	user.RegisterUserServiceServer(server, userHandler)
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
// func StreamServerInterceptorCustom() grpc.StreamServerInterceptor {
// 	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		log.Printf("StreamServerInterceptorCustom => %v", info.FullMethod)
// 		return handler(srv, stream)
// 	}
// }

// UnaryServerInterceptorCustom
// func UnaryServerInterceptorCustom() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		log.Printf("UnaryServerInterceptorCustom => %v", req)
// 		return handler(ctx, req)
// 	}
// }
