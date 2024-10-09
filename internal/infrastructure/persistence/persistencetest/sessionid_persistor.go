package persistencetest

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model"
)

type sessionIDPersistor struct {
	t  *testing.T
	db *sql.DB
}

func NewSessionIDPersistor(t *testing.T, db *sql.DB) *sessionIDPersistor {
	t.Helper()

	return &sessionIDPersistor{
		t:  t,
		db: db,
	}
}

func (p *sessionIDPersistor) Add(sessionID *model.CandhisSessionID) {
	p.t.Helper()

	_, err := p.db.Exec("INSERT INTO candhis_session (id, created_at) VALUES ($1, $2)", sessionID.ID(), sessionID.CreatedAt())
	require.NoError(p.t, err, "failed to insert session ID: %v", err)
}

func (p *sessionIDPersistor) Clear() {
	p.t.Helper()

	_, err := p.db.Exec("DELETE FROM candhis_session")
	require.NoError(p.t, err, "failed to clear candhis_session table: %v", err)
}
