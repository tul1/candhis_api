package persistencetest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
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

type ESPersistor struct {
	waveDataPersistor *waveDataPersistor
}

func NewESPersistor(t *testing.T, es *elasticsearch.Client) *ESPersistor {
	t.Helper()

	return &ESPersistor{
		waveDataPersistor: NewWaveDataPersistor(t, es),
	}
}

func (p *ESPersistor) WaveData() *waveDataPersistor {
	return p.waveDataPersistor
}

func (p *ESPersistor) Clear(ctx context.Context) {
	p.waveDataPersistor.Clear(ctx)
}
