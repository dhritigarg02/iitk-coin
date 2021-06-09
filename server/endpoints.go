package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, Web!\n"))
}

func Login(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}
		var loginUser AuthUser
		err := json.NewDecoder(r.Body).Decode(&loginUser)
		if err != nil {
			http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
			return
		}
		var password string
		row := db.QueryRow("SELECT password FROM Auth WHERE rollno = ?", loginUser.RollNo)
		result := row.Scan(&password)

		if result == sql.ErrNoRows {
			http.Error(w, "User does not exist!", http.StatusNotFound)
			return
		}
		if result != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginUser.Password))
		if err != nil {
			http.Error(w, "Invalid Password!", http.StatusUnauthorized)
			return
		}

		token, err := CreateToken(loginUser.RollNo)
		if err != nil {
			http.Error(w, "Error while generating token, please try again", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Token{token})
	}
}

func Signup(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
			return
		}

		var newUser User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
			return
		}

		notExists, err := UserExists(db, newUser.RollNo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if !(notExists) {
			http.Error(w, "User already exists!", http.StatusConflict)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 5)
		if err != nil {
			http.Error(w, "Error while Hashing Password, Please Try Again", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		newUser.Password = string(hash)
		err = Add_User(db, newUser)
		if err != nil {
			http.Error(w, "Error while creating User, Please Try Again", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = Add_auth_data(db, newUser)
		if err != nil {
			http.Error(w, "Error while creating User, Please Try Again", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Write([]byte("Signup Successful!"))
	}
}

func Secretpage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")

	_, err := VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Authorized!"))
}
