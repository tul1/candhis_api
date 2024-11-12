package candhisapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	candhisapi "github.com/tul1/candhis_api/internal/application/candhis_api"
	"github.com/tul1/candhis_api/openapi"
)

func TestWaves(t *testing.T) {
	resp := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(resp)
	api := candhisapi.NewCandhisAPI(r)

	api.Waves(ctx)

	var wavesResp openapi.Error
	err := json.Unmarshal(resp.Body.Bytes(), &wavesResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "unimplemented path", wavesResp.Message)
}
