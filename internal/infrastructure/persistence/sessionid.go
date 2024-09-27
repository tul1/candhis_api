package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/pkg/db"
)

type sessionIDRepository struct {
	db *sql.DB
}

func NewSessionIDRepository(db *sql.DB) *sessionIDRepository {
	return &sessionIDRepository{
		db: db,
	}
}

func (r *sessionIDRepository) Get(ctx context.Context) (*model.CandhisSessionID, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, created_at FROM candhis_session`)

	var sessionID model.CandhisSessionID
	err := row.Scan(&sessionID.ID, &sessionID.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no session ID found in database")
		}
		return nil, fmt.Errorf("failed to get session ID from database: %w", err)
	}

	return &sessionID, nil
}

func (r *sessionIDRepository) Update(ctx context.Context, sessionID *model.CandhisSessionID) error {
	return db.Transaction(ctx, r.db, func(tx *sql.Tx) error {
		// Update the session ID within a transaction
		query := `UPDATE candhis_session SET id = $1, created_at = $2`
		result, err := tx.ExecContext(ctx, query, sessionID.ID, sessionID.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to update session ID: %w", err)
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to check affected rows: %w", err)
		}

		if rowsAffected == 0 {
			return errors.New("no session ID found to update")
		}

		return nil
	})
}
