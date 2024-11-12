package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaves(t *testing.T) {
	openAPIClient := setupOpenAPIClient(t)

	respWaves, err := openAPIClient.WavesWithResponse(context.Background())
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, respWaves.StatusCode())

	errResp := respWaves.JSON500
	assert.Equal(t, "unimplemented path", errResp.Message)
}
