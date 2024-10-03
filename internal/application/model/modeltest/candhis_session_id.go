package modeltest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model"
)

func MustCreateCandhisSessionID(t *testing.T, id string) model.CandhisSessionID {
	t.Helper()

	sessionID, err := model.NewCandhisSessionID(id, nil)
	require.NoError(t, err, "failed to create CandhisSessionID")

	return sessionID
}
