package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tul1/candhis_api/internal/pkg/logger"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	log.Info("Test Info message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "Test Info message", logEntry["msg"])
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	log.Error("Test Error message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "error", logEntry["level"])
	assert.Equal(t, "Test Error message", logEntry["msg"])
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	log = log.WithFields(logger.Fields{
		"user_id": 123,
		"action":  "test",
	})

	log.Info("Test message with fields")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "Test message with fields", logEntry["msg"])
	assert.Equal(t, float64(123), logEntry["user_id"])
	assert.Equal(t, "test", logEntry["action"])
}

func TestLogger_ContextFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	// Create context with fields
	ctx := context.Background()
	ctx = logger.ContextWithFields(ctx, logger.Fields{
		"request_id": "abc123",
	})

	log.InfoContextf(ctx, "Test message with context fields")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "Test message with context fields", logEntry["msg"])
	assert.Equal(t, "abc123", logEntry["request_id"])
}

func TestLogger_MergeFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	log = log.WithField("session_id", "sess456")

	// Create context with fields
	ctx := context.Background()
	ctx = logger.ContextWithFields(ctx, logger.Fields{
		"user_id": 789,
	})

	log.InfoContextf(ctx, "Test message with merged fields")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "Test message with merged fields", logEntry["msg"])
	assert.Equal(t, "sess456", logEntry["session_id"])
	assert.Equal(t, float64(789), logEntry["user_id"])
}

func TestLogger_LogLevels(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)
	log.SetLevel(logrus.WarnLevel)

	log.Info("This should not be logged")
	log.Warn("This should be logged")

	var logEntries []map[string]interface{}
	decoder := json.NewDecoder(&buf)
	for decoder.More() {
		var logEntry map[string]interface{}
		err := decoder.Decode(&logEntry)
		assert.NoError(t, err)
		logEntries = append(logEntries, logEntry)
	}

	assert.Equal(t, 1, len(logEntries))
	assert.Equal(t, "warning", logEntries[0]["level"])
	assert.Equal(t, "This should be logged", logEntries[0]["msg"])
}

func TestLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	errMsg := "an error occurred"
	log.Errorf("Test Errorf: %s", errMsg)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "error", logEntry["level"])
	assert.Equal(t, "Test Errorf: an error occurred", logEntry["msg"])
}

func TestFields_Merge(t *testing.T) {
	fields1 := logger.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	fields2 := logger.Fields{
		"key2": "override",
		"key3": "value3",
	}

	mergedFields := fields1.Merge(fields2)
	expectedFields := logger.Fields{
		"key1": "value1",
		"key2": "override",
		"key3": "value3",
	}

	assert.Equal(t, expectedFields, mergedFields)
}

func TestNewFieldsFromContext(t *testing.T) {
	fields := logger.Fields{
		"user_id": 42,
	}
	ctx := context.Background()
	ctx = logger.ContextWithFields(ctx, fields)

	retrievedFields := logger.NewFieldsFromContext(ctx)
	assert.Equal(t, fields, retrievedFields)
}

func TestLogger_WithField(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	log = log.WithField("transaction_id", "tx123")
	log.Info("Test message with single field")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "tx123", logEntry["transaction_id"])
}

func TestLogger_InfoContextf(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithDefaultLogger()
	log.SetOutput(&buf)

	ctx := context.Background()
	ctx = logger.ContextWithFields(ctx, logger.Fields{
		"order_id": "order789",
	})

	log.InfoContextf(ctx, "Processing order %s", "order789")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "Processing order order789", logEntry["msg"])
	assert.Equal(t, "order789", logEntry["order_id"])
}
