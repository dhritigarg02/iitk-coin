package db

import(
	"fmt"
	"errors"
	"log"
)

var ErrInsufficientBal = errors.New("insufficient balance in sender's wallet")

func (store *DBStore) CalculateAmntRcvd(amount int, tax int) int {

	amnt := float64(amount) * float64(1.0 - float64(tax)/100)
	
	return int(amnt)
}

func (store *DBStore) GetTax(roll1 int, roll2 int) (int, error){

	var batch1, batch2 int

	statement, err := store.db.Prepare(
		`SELECT batch
		FROM User
		WHERE rollno = ?
		`)
	if err != nil {
		return 0, err
	}

	row := statement.QueryRow(roll1)
	err = row.Scan(
		&batch1,
	)
	if err != nil {
		return 0, err
	}

	row = statement.QueryRow(roll2)
	err = row.Scan(
		&batch2,
	)
	if err != nil {
		return 0, err
	}

	if batch1 == batch2 {
		return store.IntraBatchTax, nil
	} else {
		return store.InterBatchTax, nil
	}
}

func (store *DBStore) ExecTx(fn func(*Queries) error) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil { 
			return fmt.Errorf("[Tx error]: %v [Rb error]: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (store *DBStore) AddCoins(req EntryParams) error {

	err := store.ExecTx(func(q *Queries) error{

		err := q.AddEntry(req)
		if err != nil {
			return err
		}

		err = q.UpdateBalance(req)
		return err
	})

	return err
}

func (store *DBStore) TransferCoins(req TransferParams) error {

	err := store.ExecTx(func(q *Queries) error{

		statement, err := q.db.Prepare(
			`UPDATE Wallet 
			SET coins = CASE WHEN coins >= ? THEN (coins - ?) ELSE coins END
			WHERE rollno = ?
			`)
		if err != nil {
			return err
		}
		result, err := statement.Exec(req.Amount, req.Amount, req.Sender)
		if err != nil {
			return err
		}
		numRowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if numRowsAffected > 1 {
			log.Printf("[TransferCoins] [ERROR] : Duplicate rows in database!, %v", err)
			return errors.New("internal Server Error")
		} else if numRowsAffected == 0 {
			return ErrInsufficientBal
		}

		err = q.UpdateBalance(EntryParams{
			RollNo: req.Receiver,
			Amount: req.AmountRcvd,
		})
		if err != nil {
			return err
		}

		err = q.AddTransfer(req)
		if err != nil {
			return err
		}

		err = q.AddEntry(EntryParams{
			RollNo: req.Sender,
			Amount: -req.Amount,
		})
		if err != nil {
			return err
		}

		err = q.AddEntry(EntryParams{
			RollNo: req.Receiver,
			Amount: req.AmountRcvd,
		})
		return err
	})
	return err
}


