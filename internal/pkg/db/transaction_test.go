package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbpkg "github.com/tul1/candhis_api/internal/pkg/db"
)

func TestTransaction_FailureToBegin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.Background()

	mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))

	err = dbpkg.Transaction(ctx, db, func(tx *sql.Tx) error {
		return nil
	})

	assert.EqualError(t, err, "failed to begin transaction: failed to begin transaction")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_FailureDuringTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.Background()

	mock.ExpectBegin()

	mock.ExpectRollback()

	err = dbpkg.Transaction(ctx, db, func(tx *sql.Tx) error {
		return errors.New("operation failed")
	})

	assert.EqualError(t, err, "operation failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_FailureDuringRollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))

	err = dbpkg.Transaction(ctx, db, func(tx *sql.Tx) error {
		return errors.New("operation failed")
	})

	assert.EqualError(t, err, "failed to rollback transaction: rollback failed after error: operation failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.Background()

	mock.ExpectBegin()

	mock.ExpectCommit()

	err = dbpkg.Transaction(ctx, db, func(tx *sql.Tx) error {
		return nil
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
