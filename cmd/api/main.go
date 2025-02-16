package main

import (
	"log"
	"os"
	"strconv"

	"example.com/user-management/internal/infrastructure/persistence/postgres"
	"example.com/user-management/internal/interface/handler"
	"example.com/user-management/internal/usecase"
	"github.com/savsgio/atreugo/v11"
)

func main() {
	// Initialize database
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)

	// Server configuration
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if port == 0 {
		port = 8080
	}
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Initialize server
	server := atreugo.New(atreugo.Config{
		Addr: host + ":" + strconv.Itoa(port),
	})

	// User routes
	server.Path("POST", "/users", userHandler.Create)
	server.Path("GET", "/users", userHandler.List)
	server.Path("GET", "/users/:id", userHandler.Get)
	server.Path("PUT", "/users/:id", userHandler.Update)
	server.Path("DELETE", "/users/:id", userHandler.Delete)

	log.Printf("Server starting on %s:%d", host, port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
