package db

import (
	"context"
	"database/sql"
	"fmt"
)

func Transaction(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = f(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w after error: %w", txErr, err)
		}

		return err
	}

	return tx.Commit()
}
