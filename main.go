package main

import (
	"context"
	"fmt"
	"go-crud/db"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDb()
	defer db.Db.Close(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /user", createUsersHandler)
	mux.HandleFunc("GET /users", getUsersHandler)
	mux.HandleFunc("GET /user/{id}", getUserHandler)
	mux.HandleFunc("PUT /user/{id}", updateUserHandler)
	mux.HandleFunc("DELETE /user/{id}", deleteUserHandler)

	fmt.Println("Server is running on port 5000")
	err = http.ListenAndServe(":5000", mux)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
