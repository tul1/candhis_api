package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
)

func TestSessionIDRepository_Get_DatabaseError(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	mock.ExpectQuery(`SELECT id, created_at FROM candhis_session`).
		WillReturnError(errors.New("database error"))

	_, err := repo.Get(context.Background())
	assert.EqualError(t, err, "failed to get session ID from database: database error")
}

func TestSessionIDRepository_Get_NotFound(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	mock.ExpectQuery(`SELECT id, created_at FROM candhis_session`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}))

	_, err := repo.Get(context.Background())
	assert.EqualError(t, err, "no session ID found in database")
}

func TestSessionIDRepository_Get_Success(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	expectedID := "some-session-id"
	expectedCreatedAt := time.Now()

	mock.ExpectQuery(`SELECT id, created_at FROM candhis_session`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(expectedID, expectedCreatedAt))

	sessionID, err := repo.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, expectedID, sessionID.ID)
	assert.WithinDuration(t, expectedCreatedAt, sessionID.CreatedAt, time.Second)
}

func TestSessionIDRepository_Update_DatabaseError(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	sessionID := &model.CandhisSessionID{
		ID:        "some-session-id",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE candhis_session SET id = \$1, created_at = \$2`).
		WithArgs(sessionID.ID, sessionID.CreatedAt).
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	err := repo.Update(context.Background(), sessionID)
	assert.EqualError(t, err, "failed to update session ID: update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionIDRepository_Update_NoRowsAffected(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	sessionID := &model.CandhisSessionID{
		ID:        "some-session-id",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE candhis_session SET id = \$1, created_at = \$2`).
		WithArgs(sessionID.ID, sessionID.CreatedAt).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.Update(context.Background(), sessionID)
	assert.EqualError(t, err, "no session ID found to update")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionIDRepository_Update_Success(t *testing.T) {
	repo, mock := setupSessionIDSQLMock(t)

	sessionID := &model.CandhisSessionID{
		ID:        "updated-session-id",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE candhis_session SET id = \$1, created_at = \$2`).
		WithArgs(sessionID.ID, sessionID.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), sessionID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func setupSessionIDSQLMock(t *testing.T) (repository.SessionID, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := persistence.NewSessionIDRepository(db)

	return repo, mock
}
