package database

import (
	"context"
	"fmt"
)

// TxFn is a function that will be executed in a transaction.
// It takes a *Queries object that is bound to the SPECIFIC transaction.
type TxFn func(*Queries) error

// ExecTx executes a function within a database transaction
func (s *Service) ExecTx(ctx context.Context, fn TxFn) error {
	// 1. Start the transaction
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	// 2. Create a specific Query engine for this transaction
	q := New(tx)

	// 3. Run the user's function
	err = fn(q)

	// 4. Commit or Rollback
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
