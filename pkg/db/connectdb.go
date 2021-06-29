package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ConnectDB() *sql.DB {

	database, err := sql.Open("sqlite3", "./user_data.db")
	handleError(err)

	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS User (
			id INTEGER PRIMARY KEY,
			rollno INTEGER NOT NULL UNIQUE, 
			name TEXT NOT NULL, 
			batch INTEGER NOT NULL,
			isAdmin INTEGER DEFAULT 0,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
			)`)
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	statement, err = database.Prepare(
		`CREATE TABLE IF NOT EXISTS Auth (
			id INTEGER PRIMARY KEY,
			rollno INTEGER NOT NULL UNIQUE, 
			password TEXT NOT NULL)`)
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	statement, err = database.Prepare(
		`CREATE TABLE IF NOT EXISTS Wallet (
			id INTEGER PRIMARY KEY,
			rollno INTEGER NOT NULL UNIQUE, 
			coins INTEGER NOT NULL)`)
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	statement, err = database.Prepare(
		`CREATE TABLE IF NOT EXISTS Entries (
			id INTEGER PRIMARY KEY,
			rollno INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
			)`)
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	statement, err = database.Prepare(
		`CREATE TABLE IF NOT EXISTS Transfers (
			id INTEGER PRIMARY KEY,
			receiver INTEGER NOT NULL,
			sender INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			tax INTEGER,
			amountrcvd INTEGER,
			remarks TEXT,
			time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
			)`)
	handleError(err)

	_, err = statement.Exec()
	handleError(err)

	return database
}