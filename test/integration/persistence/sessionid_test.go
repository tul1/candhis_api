package persistence_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model/modeltest"
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

	// Add session ID to the database using the persistor
	sessionID := modeltest.MustCreateCandhisSessionID(t, "some-session-id")
	persistor.SessionID().Add(&sessionID)

	// Fetch the session ID from the database
	retrievedSessionID, err := sessionIDRepository.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "some-session-id", retrievedSessionID.ID())
	assert.Equal(t, sessionID.CreatedAt(), retrievedSessionID.CreatedAt())
}

func TestSessionIDRepository_Update_NotExistingCandhisSessionID(t *testing.T) {
	_, sessionIDRepository := setupSessionIDTest(t)

	// Try to update a session ID that doesn't exist
	sessionID := modeltest.MustCreateCandhisSessionID(t, "non-existing-session-id")
	err := sessionIDRepository.Update(context.Background(), &sessionID)
	assert.EqualError(t, err, "no session ID found to update")
}

func TestSessionIDRepository_Update_Success(t *testing.T) {
	persistor, sessionIDRepository := setupSessionIDTest(t)

	// Add an initial session ID to the database
	initialSessionID := modeltest.MustCreateCandhisSessionID(t, "initial-session-id")
	persistor.SessionID().Add(&initialSessionID)

	// Update the session ID with new values
	updatedSessionID := modeltest.MustCreateCandhisSessionID(t, "updated-session-id")
	err := sessionIDRepository.Update(context.Background(), &updatedSessionID)
	require.NoError(t, err)

	// Fetch the session ID from the database and assert the updated values
	retrievedSessionID, err := sessionIDRepository.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "updated-session-id", retrievedSessionID.ID())
	assert.Equal(t, updatedSessionID.CreatedAt(), retrievedSessionID.CreatedAt())
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
