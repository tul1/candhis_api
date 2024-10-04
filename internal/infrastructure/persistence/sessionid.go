package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/pkg/db"
)

type sessionIDRepository struct {
	dbConn *sql.DB
}

func NewSessionIDRepository(dbConn *sql.DB) *sessionIDRepository {
	return &sessionIDRepository{
		dbConn: dbConn,
	}
}

func (r *sessionIDRepository) Get(ctx context.Context) (*model.CandhisSessionID, error) {
	row := r.dbConn.QueryRowContext(ctx, `SELECT id, created_at FROM candhis_session`)

	var id string
	var createdAt time.Time

	err := row.Scan(&id, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no session ID found in database")
		}
		return nil, fmt.Errorf("failed to get session ID from database: %w", err)
	}

	c := createdAt.UTC()
	candhisSessionID, err := model.NewCandhisSessionID(id, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to create session ID: %w", err)
	}

	return &candhisSessionID, nil
}

func (r *sessionIDRepository) Update(ctx context.Context, sessionID *model.CandhisSessionID) error {
	return db.Transaction(ctx, r.dbConn, func(tx *sql.Tx) error {
		query := `UPDATE candhis_session SET id = $1, created_at = $2`
		result, err := tx.ExecContext(ctx, query, sessionID.ID(), sessionID.CreatedAt())
		if err != nil {
			return fmt.Errorf("failed to update session ID: %w", err)
		}

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
