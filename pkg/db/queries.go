package db

import (
	"fmt"
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
		return fmt.Errorf("[CreateUser] : %v", err)
	}
	_, err = statement.Exec(req.RollNo, req.Name, req.Batch, req.IsAdmin)
	if err != nil {
		return fmt.Errorf("[CreateUser] : %v", err)
	}

	statement, err = q.db.Prepare(
		`INSERT INTO Auth(
			rollno,
			password
		) VALUES(?, ?)
		`)
	if err != nil {
		return fmt.Errorf("[CreateUser] : %v", err)
	}
	_, err = statement.Exec(req.RollNo, req.Password)
	if err != nil {
		return fmt.Errorf("[CreateUser] : %v", err)
	}

	statement, err = q.db.Prepare(
		`INSERT INTO Wallet(
			rollno,
			coins
		) VALUES(?, ?)
		`)
	if err != nil {
		return fmt.Errorf("[CreateUser] : %v", err)
	}
	_, err = statement.Exec(req.RollNo, 0)
	if err != nil {
		return fmt.Errorf("[CreateUser] : %v", err)
	}

	return nil
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
			return false, fmt.Errorf("[UserExists] : %v", err)
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
		return fmt.Errorf("[AddEntry] : %v", err)
	}

	_, err = statement.Exec(req.RollNo, req.Amount)
	if err != nil {
		return fmt.Errorf("[AddEntry] : %v", err)
	}

	return nil
}

func (q *Queries) AddTransfer(req TransferParams) (error) {

	statement, err := q.db.Prepare(
		`INSERT INTO Transfers(
			receiver,
			sender,
			amount,
			tax,
			amountrcvd,
			remarks
		) VALUES(?, ?, ?, ?, ?, ?)
		`)
	if err != nil {
		return fmt.Errorf("[AddTransfer] : %v", err)
	}

	_, err = statement.Exec(req.Receiver, req.Sender, req.Amount, req.Tax, req.AmountRcvd, req.Remarks)
	if err != nil {
		return fmt.Errorf("[AddTransfer] : %v", err)
	}

	return nil
}

func (q *Queries) UpdateBalance(req EntryParams) (error) {

	statement, err := q.db.Prepare(
		`UPDATE Wallet 
		SET coins = coins + ?
		WHERE rollno = ?
		`)
	if err != nil {
		return fmt.Errorf("[UpdateBalance] : %v", err)
	}

	_, err = statement.Exec(req.Amount, req.RollNo)
	if err != nil {
		return fmt.Errorf("[UpdateBalance] : %v", err)
	}

	return nil
}

func (q *Queries) GetBalance(rollno int) (int, error) {

	var balance int

	statement, err := q.db.Prepare(
		`SELECT coins
		FROM Wallet
		WHERE rollno = ?
		`)
	if err != nil {
		return balance, fmt.Errorf("[GetBalance] : %v", err)
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&balance,
	)
	if err != nil {
		return balance, fmt.Errorf("[GetBalance] : %v", err)
	}

	return balance, nil
}

func (q *Queries) GetItemCost(itemid int) (int, error) {

	var cost int

	statement, err := q.db.Prepare(
		`SELECT cost
		FROM Items
		WHERE itemid = ?
		`)
	if err != nil {
		return cost, fmt.Errorf("[GetItemCost] : %v", err)
	}

	row := statement.QueryRow(itemid)
	err = row.Scan(
		&cost,
	)
	if err != nil {
		return cost, fmt.Errorf("[GetItemCost] : %v", err)
	}

	return cost, nil
}

func (q *Queries) AddRedeemReq(req Redeem) (error) {

	statement, err := q.db.Prepare(
		`INSERT INTO RedeemReqs (
			rollno,
			itemid
		) VALUES(?, ?)
		`)
	if err != nil {
		return fmt.Errorf("[AddRedeemReq] : %v", err)
	}

	_, err = statement.Exec(req.RollNo, req.ItemId)
	if err != nil {
		return fmt.Errorf("[AddRedeemReq] : %v", err)
	}

	return nil
}

func (q* Queries) GetHashedPswd(rollno int) (string, error) {

	var pswd string

	statement, err := q.db.Prepare(
		`SELECT password
		FROM Auth 
		WHERE rollno = ?
		`)
	if err != nil {
		return pswd, fmt.Errorf("[GetHashedPswd] : %v", err)
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&pswd,
	)
	if err != nil {
		return pswd, fmt.Errorf("[GetHashedPswd] : %v", err)
	}

	return pswd, nil
}

func (q* Queries) CheckAdmin(rollno int) (bool, error) {

	var isAdmin bool

	statement, err := q.db.Prepare(
		`SELECT isAdmin
		FROM User
		WHERE rollno = ?
		`)
	if err != nil {
		return false, fmt.Errorf("[CheckAdmin] : %v", err)
	}

	row := statement.QueryRow(rollno)
	err = row.Scan(
		&isAdmin,
	)
	if err != nil {
		return false, fmt.Errorf("[CheckAdmin] : %v", err)
	}

	return isAdmin, nil
}
