package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// begin a transaction : BEGIN
	tx, err := store.connPool.Begin(ctx)
	// check for errors
	if err != nil {
		return err
	}
	// create new Queries
	q := New(tx)
	// execute transations
	err = fn(q)
	// check for errors
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// return
	return tx.Commit(ctx)
}
