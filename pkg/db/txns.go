package db

import(
	"fmt"
)

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

