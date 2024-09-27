package persistencetest

import (
	"database/sql"
	"testing"

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
	_, err := p.db.Exec("INSERT INTO candhis_session (id, created_at) VALUES ($1, $2)", sessionID.ID, sessionID.CreatedAt)
	if err != nil {
		p.t.Fatalf("failed to insert session ID: %v", err)
	}
}

func (p *sessionIDPersistor) Clear() {
	_, err := p.db.Exec("DELETE FROM candhis_session")
	if err != nil {
		p.t.Fatalf("failed to clear candhis_session table: %v", err)
	}
}
