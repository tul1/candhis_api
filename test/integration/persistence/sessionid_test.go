package persistence_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model/modeltest"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence/persistencetest"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/logger"
)

func TestSessionIDStore_Get_NotFound(t *testing.T) {
	_, sessionIDStore := setupSessionIDTest(t)

	_, err := sessionIDStore.Get(context.Background())
	require.Error(t, err)

	assert.EqualError(t, err, "no session ID found in database")
}

func TestSessionIDStore_Get_Success(t *testing.T) {
	persistor, sessionIDStore := setupSessionIDTest(t)

	sessionID := modeltest.MustCreateCandhisSessionID(t, "some-session-id")
	persistor.SessionID().Add(&sessionID)

	retrievedSessionID, err := sessionIDStore.Get(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "some-session-id", retrievedSessionID.ID())
	assert.Equal(t, sessionID.CreatedAt(), retrievedSessionID.CreatedAt())
}

func TestSessionIDStore_Update_NotExistingCandhisSessionID(t *testing.T) {
	_, sessionIDStore := setupSessionIDTest(t)

	sessionID := modeltest.MustCreateCandhisSessionID(t, "non-existing-session-id")
	err := sessionIDStore.Update(context.Background(), sessionID)
	assert.EqualError(t, err, "no session ID found to update")
}

func TestSessionIDStore_Update_Success(t *testing.T) {
	persistor, sessionIDStore := setupSessionIDTest(t)

	initialSessionID := modeltest.MustCreateCandhisSessionID(t, "initial-session-id")
	persistor.SessionID().Add(&initialSessionID)

	updatedSessionID := modeltest.MustCreateCandhisSessionID(t, "updated-session-id")
	err := sessionIDStore.Update(context.Background(), updatedSessionID)
	require.NoError(t, err)

	retrievedSessionID, err := sessionIDStore.Get(context.Background())
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

	dbName := os.Getenv("DATABASE_NAME")
	require.NotEmpty(t, dbName)

	password := os.Getenv("DATABASE_PASSWORD")
	require.NotEmpty(t, password)

	dbConn, err := db.NewDBConnection(user, password, host, port, dbName, db.DefaultDBConnector, logger.NewWithDefaultLogger())
	require.NoError(t, err, "failed to initialize database connection")

	sessionIDStore := persistence.NewSessionID(dbConn.DB)
	persistor := persistencetest.NewPersistor(t, dbConn.DB)

	t.Cleanup(func() { persistor.Clear() })

	return persistor, sessionIDStore
}
