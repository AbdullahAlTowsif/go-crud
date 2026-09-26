package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var users = []User{
	{
		Id:    1,
		Name:  "Towsif",
		Age:   23,
		Email: "towsif@gmail.com",
	},
	{
		Id:    2,
		Name:  "Abdullah",
		Age:   24,
		Email: "abdullah@gmail.com",
	},
}

var db *pgx.Conn

func connectDb() {
	var err error
	connStr := os.Getenv("DB_URL")

	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected Successfully!")
}

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectDb()
	defer db.Close(context.Background())

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to go server!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is up and running!")
}

func createUsersHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}
	// newUser.Id = len(users) + 1
	// users = append(users, newUser)

	query := `
		insert into users (name, age, email)
		values ($1, $2, $3)
		returning id
	`
	err = db.QueryRow(context.Background(), query, newUser.Name, newUser.Age, newUser.Email).Scan(&newUser.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not create user!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		select id, name, age, email from users
	`

	rows, err := db.Query(context.Background(), query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not get users")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Could not get users")
			return
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not read users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// users, _ := json.Marshal(users) // converts the data into json after saving into the memory
	// w.Write(users)

	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var user User

	query := `
		select id, name, age, email from users where id = $1
	`

	err = db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Name, &user.Age, &user.Email)

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User Not Found")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not get users")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	var updatedUser User

	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	query := `
		update users
		set name = $1, age = $2, email = $3
		where id = $4
		returning id, name, age, email
	`

	err = db.QueryRow(context.Background(), query, updatedUser.Name, updatedUser.Age, updatedUser.Email, id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Age, &updatedUser.Email)

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User Not Found")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not update user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user id")
		return
	}

	query := `
		delete from users
		where id = $1
	`

	var cmdTag pgconn.CommandTag

	cmdTag, err = db.Exec(context.Background(), query, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not delete user")
		return
	}

	if cmdTag.RowsAffected() != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "User Deleted Successfully!")
}
