package e2e_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/openapi"
)

func setupOpenAPIClient(t *testing.T) *openapi.ClientWithResponses {
	t.Helper()

	publicURL := os.Getenv("API_PUBLIC_URL")
	require.NotEmpty(t, publicURL, "API_PUBLIC_URL should not be empty")

	openAPIClient, err := openapi.NewClientWithResponses(publicURL)
	require.NoError(t, err)

	return openAPIClient
}
