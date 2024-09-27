package persistence_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence/persistencetest"
	"github.com/tul1/candhis_api/internal/pkg/db"
)

func TestSessionIDRepository_Get_NotFound(t *testing.T) {
	_, sessionIDRepository := setupSessionIDTest(t)

	_, err := sessionIDRepository.Get(context.Background())
	require.Error(t, err)

	assert.EqualError(t, err, "no session ID found in database")
}

func TestSessionIDRepository_Get_Success(t *testing.T) {
	persistor, sessionIDRepository := setupSessionIDTest(t)

	createdAt := time.Now().Truncate(time.Second).UTC()
	sessionID := &model.CandhisSessionID{
		ID:        "test-session-id",
		CreatedAt: createdAt,
	}

	// Add session ID to the database using the persistor
	persistor.SessionID().Add(sessionID)

	// Fetch the session ID from the database
	retrievedSessionID, err := sessionIDRepository.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, sessionID.ID, retrievedSessionID.ID)
	assert.Equal(t, createdAt, retrievedSessionID.CreatedAt)
}

func TestSessionIDRepository_Update_NoRowsAffected(t *testing.T) {
	_, sessionIDRepository := setupSessionIDTest(t)

	sessionID := &model.CandhisSessionID{
		ID:        "non-existing-session-id",
		CreatedAt: time.Now(),
	}

	// Try to update a session ID that doesn't exist
	err := sessionIDRepository.Update(context.Background(), sessionID)
	assert.EqualError(t, err, "no session ID found to update")
}

func TestSessionIDRepository_Update_Success(t *testing.T) {
	persistor, sessionIDRepository := setupSessionIDTest(t)

	// Add an initial session ID to the database
	initialSessionID := &model.CandhisSessionID{
		ID:        "initial-session-id",
		CreatedAt: time.Now(),
	}

	persistor.SessionID().Add(initialSessionID)

	// Update the session ID with new values
	updatedSessionIDCreateAt := time.Now().Add(time.Hour).Truncate(time.Second).UTC()
	updatedSessionID := &model.CandhisSessionID{
		ID:        "updated-session-id",
		CreatedAt: updatedSessionIDCreateAt,
	}

	err := sessionIDRepository.Update(context.Background(), updatedSessionID)
	require.NoError(t, err)

	// Fetch the session ID from the database and assert the updated values
	retrievedSessionID, err := sessionIDRepository.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, updatedSessionID.ID, retrievedSessionID.ID)
	assert.Equal(t, updatedSessionIDCreateAt, retrievedSessionID.CreatedAt)
}

func setupSessionIDTest(t *testing.T) (*persistencetest.Persistor, repository.SessionID) {
	t.Helper()

	host := os.Getenv("DATABASE_HOST")
	require.NotEmpty(t, host)

	port := os.Getenv("DATABASE_PORT")
	require.NotEmpty(t, port)

	user := os.Getenv("DATABASE_USER")
	require.NotEmpty(t, user)

	dbname := os.Getenv("DATABASE_NAME")
	require.NotEmpty(t, dbname)

	password := os.Getenv("DATABASE_PASSWORD")
	require.NotEmpty(t, password)

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	db, err := db.NewDatabaseConnection(dsn)
	require.NoError(t, err, "failed to initialize database connection")

	sessionIDRepository := persistence.NewSessionIDRepository(db)
	persistor := persistencetest.NewPersistor(t, db)

	t.Cleanup(func() { persistor.Clear() })

	return persistor, sessionIDRepository
}
