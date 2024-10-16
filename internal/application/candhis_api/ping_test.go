package candhisapi_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	candhisapi "github.com/tul1/candhis_api/internal/application/candhis_api"
	"github.com/tul1/candhis_api/openapi"
)

func TestPing(t *testing.T) {
	resp := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(resp)
	api := candhisapi.NewCandhisAPI(r)

	api.Ping(ctx)

	var pong openapi.Pong
	err := json.Unmarshal(resp.Body.Bytes(), &pong)
	require.NoError(t, err)
	assert.Equal(t, openapi.Pong{Message: "pong"}, pong)
}
