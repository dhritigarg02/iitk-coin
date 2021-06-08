package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dhritigarg02/iitk-coin/server"
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

	http.HandleFunc("/login", server.Login(database))
	http.HandleFunc("/signup", server.Signup(database))

	fmt.Println("Starting server at port 8080....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
