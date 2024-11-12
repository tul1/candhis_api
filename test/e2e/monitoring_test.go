package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	openAPIClient := setupOpenAPIClient(t)

	respPing, err := openAPIClient.PingWithResponse(context.Background())
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, respPing.StatusCode())

	pong := respPing.JSON200
	assert.Equal(t, "pong", pong.Message)
}
