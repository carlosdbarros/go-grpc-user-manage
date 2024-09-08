package main

import (
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	"os"
	"path/filepath"
)

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	envDir := filepath.Join(baseDir, "cmd", "grpc-server")
	config, err := configs.LoadConfig(envDir)
	if err != nil {
		panic(err)
	}
	println("Configuration loaded")
	println("DBDriver: ", config.DBDriver)
}
