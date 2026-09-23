package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /user", createUsersHandler)
	mux.HandleFunc("GET /users", getUsersHandler)
	mux.HandleFunc("GET /user/{id}", getUserHandler)
	mux.HandleFunc("PUT /user/{id}", updateUserHandler)
	mux.HandleFunc("DELETE /user/{id}", deleteUserHandler)

	fmt.Println("Server is running on port 5000")
	err := http.ListenAndServe(":5000", mux)

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
	newUser.Id = len(users) + 1
	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
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

	for _, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User Not Found")
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

	for idx, user := range users {
		if user.Id == id {
			updatedUser.Id = id
			users[idx] = updatedUser

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User Not Found")
}


func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	//
}