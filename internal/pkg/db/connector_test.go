package db_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tul1/candhis_api/internal/pkg/db"
)

func TestNewDatabaseConnection_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectPing()

	mockConnector := func(dbURL string) (*sql.DB, error) {
		return sqlDB, nil
	}

	dbConn, err := db.NewDBConnection("user", "password", "host", "port", "dbname", mockConnector, nil)
	require.NoError(t, err)

	err = dbConn.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expected all mock expectations to be met")
}

func TestNewDatabaseConnection_InvalidCredentials(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockConnector := func(dbURL string) (*sql.DB, error) {
		return nil, errors.New("failed to init db")
	}

	_, err = db.NewDBConnection("invalid_user", "invalid_pass", "host", "port", "dbname", mockConnector, nil)
	assert.EqualError(t, err, "failed to connect to PostgreSQL: failed to init db")
}

func TestCloseWithLog(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)

	mockConnector := func(dbURL string) (*sql.DB, error) {
		return sqlDB, nil
	}

	mockLogger := &MockLogger{}
	dbConn, err := db.NewDBConnection("invalid_user", "invalid_pass", "host", "port", "dbname", mockConnector, mockLogger)
	require.NoError(t, err)

	dbConn.CloseWithLog()

	require.Len(t, mockLogger.messages, 1, "Expected one log message")
	assert.Contains(t, mockLogger.messages[0], "Failed closing database")
}

type MockLogger struct {
	messages []string
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}
