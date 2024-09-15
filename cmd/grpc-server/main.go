package main

import (
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	grpcInfra "github.com/carlosdbarros/go-grpc-user-manage/internal/infra/grpc"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/pb/permission"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/pb/user"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	db, err := configs.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create route handlers
	userHandler := grpcInfra.NewUserHandler(database.NewUserDB(db))
	permHandler := grpcInfra.NewPermissionHandler(database.NewPermissionDB(db))

	// create a gRPC server instance
	server := grpc.NewServer()

	// register the service intances with the grpc server
	user.RegisterUserServiceServer(server, userHandler)
	permission.RegisterPermissionServiceServer(server, permHandler)
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
