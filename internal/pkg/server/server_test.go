package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/pkg/server"
)

func TestNewGinServer(t *testing.T) {
	s, err := server.NewGinServer(logrus.New(), "http://localhost", 8080)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNoRouteHandler(t *testing.T) {
	s, err := server.NewGinServer(logrus.New(), "http://localhost", 8080)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/invalid", http.NoBody)
	w := httptest.NewRecorder()

	s.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "\"invalid API path\"", w.Body.String())
}

type loggerRecorder struct {
	messages []string
}

func (m *loggerRecorder) Write(p []byte) (n int, err error) {
	m.messages = append(m.messages, string(p))
	return len(p), nil
}

func TestLogRequestMiddleware(t *testing.T) {
	recorder := &loggerRecorder{}
	log := logrus.New()
	log.SetOutput(recorder)
	s, err := server.NewGinServer(log, "http://localhost", 8080)
	require.NoError(t, err)

	s.GetRouter().POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.GetRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String())

	require.Len(t, recorder.messages, 2)
	assert.Contains(t, recorder.messages[0],
		`level=info msg="request received" body="{\"key\": \"value\"}" method=POST path=/test path_rule=/test query=`)
	assert.Contains(t, recorder.messages[1],
		`level=info msg="request handled" body="{\"key\": \"value\"}" end=`)
}
