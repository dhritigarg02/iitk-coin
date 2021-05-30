package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Info struct{
	rollno int
	name string
}

func add_data(database *sql.DB, row Info) {

	statement, _ := database.Prepare("INSERT INTO User(rollno, name) VALUES(?, ?)")
	statement.Exec(row.rollno, row.name)
}

func main() {

	database, _ := sql.Open("sqlite3", "./user_data.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS User (rollno INTEGER PRIMARY KEY, name TEXT)")
	statement.Exec()

	add_data(database, Info{190289, "Dhriti"})
	add_data(database, Info{190458, "Harshit"})
}
