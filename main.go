package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dhritigarg02/iitk-coin/server"
	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func connectDB() *sql.DB {

	database, err := sql.Open("sqlite3", "./user_data.db")
	handleError(err)

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS User (rollno INTEGER PRIMARY KEY, name TEXT, batch INTEGER)")
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Auth (rollno INTEGER PRIMARY KEY, password TEXT)")
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	return database
}
func main() {

	database := connectDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HelloHandler)
	mux.HandleFunc("/login", server.Login(database))
	mux.HandleFunc("/signup", server.Signup(database))
	mux.HandleFunc("/secretpage", server.Secretpage)

	log.Println("[MAIN] [INFO] Starting server at port 8080....")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, mux)))
}
