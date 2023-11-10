package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// create new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// begin a transaction : BEGIN
	tx, err := store.db.BeginTx(ctx, nil)
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// return
	return tx.Commit()
}

// transfer tx params
type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

// transfer result
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	Amount      int64    `json:"amount"`
}

// specific transation
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}
		// update balence
		if arg.FromAccountID.Time() < arg.ToAccountID.Time() {
			result.FromAccount, result.ToAccount, err = moveMoney(ctx, q, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			}, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = moveMoney(ctx, q, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			}, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func moveMoney(ctx context.Context, q *Queries, AddAccBa1, AddAccBa2 AddAccountBalanceParams) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccBa1)
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccBa2)
	if err != nil {
		return
	}

	return
}
