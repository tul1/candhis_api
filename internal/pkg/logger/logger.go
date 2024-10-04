package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// Fields represents a map of custom fields to attach to log entries.
type Fields map[string]interface{}

// Logger wraps logrus.Logger to provide structured logging with additional features.
type Logger struct {
	logger *logrus.Logger
	fields Fields
}

// NewWithDefaultLogger creates a new Logger with the default logrus logger.
func NewWithDefaultLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return &Logger{
		logger: logger,
		fields: make(Fields),
	}
}

// SetOutput sets the output destination for the logger (e.g., file, stdout).
func (l *Logger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

// SetLevel sets the logging level (e.g., Info, Warn, Error).
func (l *Logger) SetLevel(level logrus.Level) {
	l.logger.SetLevel(level)
}

// WithField adds a single field to the logger and returns a new instance.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newFields := l.fields.Copy()
	newFields[key] = value
	return &Logger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithFields adds multiple fields to the logger and returns a new instance.
func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		logger: l.logger,
		fields: l.fields.Merge(fields),
	}
}

// Info logs a message at the Info level with optional fields.
func (l *Logger) Info(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Info(msg)
}

// Infof logs a formatted message at the Info level with optional fields.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Infof(format, args...)
}

// Warn logs a message at the Warn level with optional fields.
func (l *Logger) Warn(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Warn(msg)
}

// Warnf logs a formatted message at the Warn level with optional fields.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Warnf(format, args...)
}

// Error logs a message at the Error level with optional fields.
func (l *Logger) Error(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Error(msg)
}

// Errorf logs a formatted message at the Error level with optional fields.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Errorf(format, args...)
}

// logWithContext is an internal helper to log a message with fields from context.
func (l *Logger) logWithContext(ctx context.Context, level logrus.Level, msg string, args ...interface{}) {
	mergedFields := l.fields.Merge(NewFieldsFromContext(ctx))
	l.logger.WithFields(logrus.Fields(mergedFields)).Log(level, fmt.Sprintf(msg, args...))
}

// InfoContextf logs a formatted message at the Info level using fields from context.
func (l *Logger) InfoContextf(ctx context.Context, format string, args ...interface{}) {
	l.logWithContext(ctx, logrus.InfoLevel, format, args...)
}

// WarnContextf logs a formatted message at the Warn level using fields from context.
func (l *Logger) WarnContextf(ctx context.Context, format string, args ...interface{}) {
	l.logWithContext(ctx, logrus.WarnLevel, format, args...)
}

// ErrorContextf logs a formatted message at the Error level using fields from context.
func (l *Logger) ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	l.logWithContext(ctx, logrus.ErrorLevel, format, args...)
}

// Copy creates a deep copy of the Fields map.
func (f Fields) Copy() Fields {
	newFields := make(Fields, len(f))
	for k, v := range f {
		newFields[k] = v
	}
	return newFields
}

// Merge merges another Fields map into the current Fields map.
func (f Fields) Merge(other Fields) Fields {
	merged := f.Copy()
	for k, v := range other {
		merged[k] = v
	}
	return merged
}

// NewFieldsFromContext retrieves the Fields stored in the context.
// If no fields are present, it returns an empty Fields map.
func NewFieldsFromContext(ctx context.Context) Fields {
	if fields, ok := ctx.Value("logFields").(Fields); ok {
		return fields
	}
	return Fields{}
}

// ContextWithFields adds custom fields to the context and returns the updated context.
func ContextWithFields(ctx context.Context, fields Fields) context.Context {
	currentFields := NewFieldsFromContext(ctx)
	mergedFields := currentFields.Merge(fields)
	return context.WithValue(ctx, "logFields", mergedFields)
}
