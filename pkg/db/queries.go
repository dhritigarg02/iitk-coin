package db

import (
	"database/sql"
)

func (q *Queries) CreateUser(req User) (error) {

	statement, err := q.db.Prepare(
		`INSERT INTO User(
			rollno,
			name,
			batch,
			isAdmin
		) VALUES(?, ?, ?, ?)
		`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(req.RollNo, req.Name, req.Batch, req.IsAdmin)
	if err != nil {
		return err
	}

	statement, err = q.db.Prepare(
		`INSERT INTO Auth(
			rollno,
			password
		) VALUES(?, ?)
		`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(req.RollNo, req.Password)
	if err != nil {
		return err
	}

	statement, err = q.db.Prepare(
		`INSERT INTO Wallet(
			rollno,
			coins
		) VALUES(?, ?)
		`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(req.RollNo, 0)
	return err
}

func (q *Queries) UserExists(rollno int) (bool, error) {

	var name string
	row := q.db.QueryRow(
		`SELECT name FROM User
		WHERE rollno = ?`, 
		rollno)
	switch err := row.Scan(&name); err {
		case sql.ErrNoRows:
			return false, nil
		case nil:
			return true, nil
		default:
			return false, err
	}
}

func (q *Queries) AddEntry(req EntryParams) (error) {

	statement, err := q.db.Prepare(
		`INSERT INTO Entries(
			rollno,
			amount
		) VALUES(?, ?)
		`)
	if err != nil {
		return err
	}

	_, err = statement.Exec(req.RollNo, req.Amount)
	return err
}

func (q *Queries) AddTransfer(req TransferParams) (error) {

	statement, err := q.db.Prepare(
		`INSERT INTO Transfers(
			receiver,
			sender,
			amount,
			remarks
		) VALUES(?, ?, ?, ?, ?)
		`)
	if err != nil {
		return err
	}

	_, err = statement.Exec(req.Receiver, req.Sender, req.Amount, req.Remarks)
	return err
}

func (q *Queries) UpdateBalance(req EntryParams) (error) {

	statement, err := q.db.Prepare(
		`UPDATE Wallet 
		SET coins = coins + ?
		WHERE rollno = ?
		`)
	if err != nil {
		return err
	}

	_, err = statement.Exec(req.Amount, req.RollNo)
	return err
}

func (q *Queries) GetBalance(rollno int) (int, error) {

	var balance int

	statement, err := q.db.Prepare(
		`SELECT coins
		FROM Wallet
		WHERE rollno = ?
		`)
	if err != nil {
		return balance, err
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&balance,
	)
	return balance, err
}

func (q* Queries) GetHashedPswd(rollno int) (string, error) {

	var pswd string

	statement, err := q.db.Prepare(
		`SELECT password
		FROM Auth 
		WHERE rollno = ?
		`)
	if err != nil {
		return pswd, err
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&pswd,
	)
	return pswd, err
}

func (q* Queries) CheckAdmin(rollno int) (bool, error) {

	var isAdmin bool

	statement, err := q.db.Prepare(
		`SELECT isAdmin
		FROM User
		WHERE rollno = ?
		`)
	if err != nil {
		return false, err
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&isAdmin,
	)
	return isAdmin, err
}
