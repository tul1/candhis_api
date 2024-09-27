package persistencetest

import (
	"database/sql"
	"testing"
)

type Persistor struct {
	sessionIDPersistor *sessionIDPersistor
}

func NewPersistor(t *testing.T, db *sql.DB) *Persistor {
	t.Helper()

	return &Persistor{
		sessionIDPersistor: NewSessionIDPersistor(t, db),
	}
}

func (p *Persistor) SessionID() *sessionIDPersistor {
	return p.sessionIDPersistor
}

func (p *Persistor) Clear() {
	p.sessionIDPersistor.Clear()
}
