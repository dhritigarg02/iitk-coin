package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/dhritigarg02/iitk-coin/pkg/auth"
	"github.com/dhritigarg02/iitk-coin/pkg/db"
)

func (server *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Web!\n"))
}

// Login logs the user in by comparing the password input
// by the user to the hashed password stored in the database.
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

// Signup creates a new user account.
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

// GetBalance retrieves the wallet balance of the logged-in user.
func (server *Server) GetBalance(w http.ResponseWriter, r *http.Request, rollno int) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	balance, err := server.DBstore.GetBalance(rollno)

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

// RewardCoins generates and adds coins to the system
// by rewarding them to the users. 
func (server *Server) RewardCoins(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var reward db.EntryParams
	err := json.NewDecoder(r.Body).Decode(&reward)
	if err != nil {
		http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
		return
	}

	if reward.RollNo == 0 || reward.Amount == 0 {
		http.Error(w, "Some fields are missing!", http.StatusBadRequest)
		return
	}

	exists, err := server.DBstore.UserExists(reward.RollNo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[RewardCoins] [ERROR] : %v\n", err)
		return
	}
	if !exists {
		http.Error(w, "User does not exist!", http.StatusNotFound)
		return
	}

	isAdmin, err := server.DBstore.CheckAdmin(reward.RollNo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[RewardCoins] [ERROR] : %v\n", err)
		return
	}
	if isAdmin {
		http.Error(w, "Admins cannot reward coins to themselves!!!", http.StatusBadRequest)
		return
	}

	err = server.DBstore.AddCoins(reward)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[RewardCoins] [ERROR] : %v\n", err)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d coins rewarded to %d", reward.Amount, reward.RollNo)))
}

// TransferCoins allows transfer of coins between two users,
// a certain percentage of coin involved is destroyed in the form of taxes.
func (server *Server) TransferCoins(w http.ResponseWriter, r *http.Request, rollno int) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var transferReq db.TransferParams
	err := json.NewDecoder(r.Body).Decode(&transferReq)
	if err != nil {
		http.Error(w, "Invalid Json provided", http.StatusUnprocessableEntity)
		return
	}

	if transferReq.Receiver == 0 || transferReq.Amount == 0 {
		http.Error(w, "Some fields are missing!", http.StatusBadRequest)
		return
	}

	transferReq.Sender = rollno
	
	exists, err := server.DBstore.UserExists(transferReq.Receiver)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[TransferCoins] [ERROR] : %v\n", err)
		return
	}
	if !exists {
		http.Error(w, "Receiver does not exist!", http.StatusNotFound)
		return
	}

	transferReq.Tax, err = server.DBstore.GetTax(transferReq.Sender, transferReq.Receiver)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[TransferCoins] [ERROR] : %v\n", err)
		return
	}
	transferReq.AmountRcvd = server.DBstore.CalculateAmntRcvd(transferReq.Amount, transferReq.Tax)
	
	err = server.DBstore.TransferCoins(transferReq)
	if err == db.ErrInsufficientBal {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("[TransferCoins] [ERROR] : %v\n", err)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d coins transferred by %d to %d\nAmount received by %d after tax deduction of %d%% is %d",
								transferReq.Amount, transferReq.Sender, transferReq.Receiver,
								transferReq.Receiver, transferReq.Tax, transferReq.AmountRcvd)))
}

