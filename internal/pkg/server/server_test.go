package server_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/pkg/logger"
	"github.com/tul1/candhis_api/internal/pkg/server"
)

type MockLogger struct {
	messages []string
	fields   logger.Fields
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func (m *MockLogger) WithFields(fields logger.Fields) server.Logger {
	m.fields = fields
	return m
}

func (m *MockLogger) Info(msg string) {
	m.messages = append(m.messages, msg)
}

func TestNewGinServer(t *testing.T) {
	mockLogger := &MockLogger{}
	s, err := server.NewGinServer(mockLogger, "http://localhost", 8080)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNoRouteHandler(t *testing.T) {
	mockLogger := &MockLogger{}
	s, err := server.NewGinServer(mockLogger, "http://localhost", 8080)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/invalid", http.NoBody)
	w := httptest.NewRecorder()

	s.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "\"invalid API path\"", w.Body.String())
}

func TestLogRequestMiddleware(t *testing.T) {
	mockLogger := &MockLogger{}
	s, err := server.NewGinServer(mockLogger, "http://localhost", 8080)
	require.NoError(t, err)

	s.GetRouter().POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.GetRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, "{\"message\":\"success\"}", w.Body.String())

	require.Len(t, mockLogger.messages, 2)
	assert.Equal(t, mockLogger.messages[0], "request received")
	assert.Equal(t, mockLogger.messages[1], "request handled")
}
