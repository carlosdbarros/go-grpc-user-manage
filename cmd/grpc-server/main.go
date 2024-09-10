package main

import (
	"database/sql"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	grpc2 "github.com/carlosdbarros/go-grpc-user-manage/internal/infra/grpc"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/pb"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := database.NewUserDBRepository(db)
	userHandler := grpc2.NewUserHandler(userRepo)

	// Create a new gRPC server
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userHandler)
	reflection.Register(server)

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	// Serve gRPC server
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	userStmt, err := db.Prepare("create table if not exists users (id text, name text, email text, password text)")
	if err != nil {
		panic(err)
	}
	_, err = userStmt.Exec()
	if err != nil {
		panic(err)
	}
	return db, nil
}
