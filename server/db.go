package server

import (
	"database/sql"
)

func Add_auth_data(db *sql.DB, row User) error {

	statement, err := db.Prepare("INSERT INTO Auth(rollno, password) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(row.RollNo, row.Password)
	return err
}

func Add_User(db *sql.DB, row User) error {

	statement, err := db.Prepare("INSERT INTO User(rollno, name, batch) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	statement.Exec(row.RollNo, row.Name, row.Batch)
	return err
}

func UserExists(db *sql.DB, rollno int) (bool, error) {
	var name string
	row := db.QueryRow("SELECT rollno FROM User WHERE rollno = ?", rollno)
	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		return true, nil
	case nil:
		return false, nil
	default:
		return false, err
	}
}
