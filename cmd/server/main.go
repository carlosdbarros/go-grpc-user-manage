package main

import (
	"database/sql"
	"encoding/json"
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/entity"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	//baseDir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//envDir := filepath.Join(baseDir, "cmd", "server")
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	initDB(db)

	userRepo := database.NewUserDBRepository(db)
	userHandler := NewUserHandler(userRepo)

	http.HandleFunc("/users", userHandler.CreateUserHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type UserHandler struct {
	repo database.UserRepository
}

func NewUserHandler(repo database.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.repo.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func initDB(db *sql.DB) {
	// Create users table if not exists
	stmt, err := db.Prepare("create table if not exists users (id text, name text, email text, password text)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}
