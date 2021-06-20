package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dhritigarg02/iitk-coin/pkg/auth"
	"github.com/dhritigarg02/iitk-coin/pkg/db"
)

func (server *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Web!\n"))
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	var loginUser db.AuthUser
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
		return
	}

	if loginUser.RollNo == 0 || loginUser.Password == "" {
		http.Error(w, "Some fields are missing!", http.StatusBadRequest)
		return
	}

	hashedpswd, err := server.DBstore.GetHashedPswd(loginUser.RollNo)
	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist!", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[Login] [ERROR] : %v\n", err)
		return
	}

	err = auth.CheckPswd(loginUser.Password, hashedpswd)
	if err != nil {
		http.Error(w, "Invalid Password!", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(loginUser.RollNo)
	if err != nil {
		http.Error(w, "Error while generating token, please try again", http.StatusInternalServerError)
		log.Printf("[Login] [ERROR] : %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}

func (server *Server) Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var newUser db.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
		return
	}

	if newUser.Name == "" || newUser.Password == "" || newUser.RollNo == 0 || newUser.Batch == 0 {
		http.Error(w, "Some fields are missing!", http.StatusBadRequest)
		return
	}

	exists, err := server.DBstore.UserExists(newUser.RollNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[Signup] [ERROR] : %v\n", err)
		return
	}
	if exists {
		http.Error(w, "User already exists!", http.StatusConflict)
		return
	}

	hashedpswd, err := auth.HashPswd(newUser.Password)
	if err != nil {
		http.Error(w, "Error while Hashing Password, Please Try Again", http.StatusInternalServerError)
		log.Printf("[Signup] [ERROR] : %v\n", err)
		return
	}
	newUser.Password = hashedpswd

	err = server.DBstore.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Error while creating User, Please Try Again", http.StatusInternalServerError)
		log.Printf("[Signup] [ERROR] : %v\n", err)
		return
	}
	w.Write([]byte("Signup Successful!"))
}

func (server *Server) Secretpage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Authorized!"))
}

func (server *Server) GetBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var user db.RollNo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
		return
	}

	if user.RollNo == 0 {
		http.Error(w, "Some fields are missing!", http.StatusBadRequest)
		return
	}

	balance, err := server.DBstore.GetBalance(user.RollNo)

	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist!", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[GetBalance] [ERROR] : %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"coins":balance})
}
