package main

import (
	"encoding/json"
	"github.com/carlosdbarros/go-grpc-user-manage/configs"
	userDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/infra/database"
	// "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"log"
	"net/http"
)

func main() {
	db, err := configs.InitSqliteInMemory()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// Create a new NewRouter and register the middleware
	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	// Create route handlers and bind them to the router
	userRepo := database.NewUserDB(db)
	userHandler := NewHttpUserHandler(userRepo)
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.FindAllUsers)

	// Start the server
	addr := "0.0.0.0:8080"
	log.Println("🚀 Server listening on: ", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type HttpUserHandler struct {
	repo userDomain.UserRepository
}

func NewHttpUserHandler(repo userDomain.UserRepository) *HttpUserHandler {
	return &HttpUserHandler{repo: repo}
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HttpUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := userDomain.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err = h.repo.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//log.Printf("Successfully created user: %v", user)
}

func (h *HttpUserHandler) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.FindAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
