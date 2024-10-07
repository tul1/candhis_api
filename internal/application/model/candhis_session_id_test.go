package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/model"
)

func TestNewCandhisSessionIDSuccess(t *testing.T) {
	testCases := map[string]struct {
		id        string
		createdAt *time.Time
	}{
		"valid session ID with UTC createdAt": {
			id:        "validSessionID123",
			createdAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
		},
		"valid session ID without createdAt": {
			id:        "validSessionID456",
			createdAt: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sessionID, err := model.NewCandhisSessionID(tc.id, tc.createdAt)
			require.NoError(t, err)
			assert.Equal(t, tc.id, sessionID.ID())

			if tc.createdAt != nil {
				assert.Equal(t, *tc.createdAt, sessionID.CreatedAt())
			} else {
				assert.NotNil(t, sessionID.CreatedAt())
			}
		})
	}
}

func TestNewCandhisSessionIDFailure(t *testing.T) {
	now := time.Now()
	nonUTC := now.In(time.FixedZone("Non-UTC", 3600)) // Non-UTC time zone

	testCases := map[string]struct {
		id        string
		createdAt *time.Time
		errMsg    string
	}{
		"invalid PHPSESSID ID": {
			id:        "PHPSESSID=invalid",
			createdAt: nil,
			errMsg:    "invalid session ID: contains PHPSESSID prefix",
		},
		"empty session ID": {
			id:        "",
			createdAt: nil,
			errMsg:    "invalid session ID: cannot be empty",
		},
		"non-UTC createdAt": {
			id:        "validSessionID789",
			createdAt: &nonUTC,
			errMsg:    "invalid createdAt: must be in UTC format",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sessionID, err := model.NewCandhisSessionID(tc.id, tc.createdAt)
			assert.EqualError(t, err, tc.errMsg)
			assert.Equal(t, model.CandhisSessionID{}, sessionID)
		})
	}
}

func TestCandhisSessionIDGetters(t *testing.T) {
	createdAt := time.Now().UTC()
	sessionID, err := model.NewCandhisSessionID("validSessionID123", &createdAt)
	require.NoError(t, err)

	assert.Equal(t, "validSessionID123", sessionID.ID())
	assert.Equal(t, "PHPSESSID=validSessionID123", sessionID.PHPSESSID())
	assert.Equal(t, createdAt, sessionID.CreatedAt())

	sessionIDWithDefaultTime, err := model.NewCandhisSessionID("validSessionID456", nil)
	require.NoError(t, err)
	assert.NotNil(t, sessionIDWithDefaultTime.CreatedAt())
}
