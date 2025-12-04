package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// LogLevel represents log levels
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// Configuration for logger
type Config struct {
	Level      LogLevel `json:"level" env:"LOG_LEVEL" envDefault:"info"`
	Format     string   `json:"format" env:"LOG_FORMAT" envDefault:"json"`   // json or text
	Output     string   `json:"output" env:"LOG_OUTPUT" envDefault:"stdout"` // stdout, stderr, or file path
	EnableFile bool     `json:"enable_file" env:"LOG_ENABLE_FILE" envDefault:"false"`
	FilePath   string   `json:"file_path" env:"LOG_FILE_PATH" envDefault:"logs/app.log"`
}

// ContextKey represents context keys for logging
type ContextKey string

const (
	CorrelationIDKey ContextKey = "correlation_id"
	UserIDKey        ContextKey = "user_id"
	RequestIDKey     ContextKey = "request_id"
	OperationKey     ContextKey = "operation"
	ServiceKey       ContextKey = "service"
)

// Global logger instance
var defaultLogger *Logger

// Initialize sets up the global logger
func Initialize(config Config) *Logger {
	logger := logrus.New()

	// Set log level
	switch config.Level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set formatter
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output
	var outputs []io.Writer

	switch config.Output {
	case "stdout":
		outputs = append(outputs, os.Stdout)
	case "stderr":
		outputs = append(outputs, os.Stderr)
	default:
		// Assume it's a file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.WithError(err).Fatal("Failed to open log file")
		}
		outputs = append(outputs, file)
	}

	// Add file output if enabled
	if config.EnableFile {
		// Create logs directory if it doesn't exist
		if err := os.MkdirAll("logs", 0755); err != nil {
			logger.WithError(err).Fatal("Failed to create logs directory")
		}

		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.WithError(err).Fatal("Failed to open log file")
		}
		outputs = append(outputs, file)
	}

	// Set multiple outputs if needed
	if len(outputs) > 1 {
		logger.SetOutput(io.MultiWriter(outputs...))
	} else if len(outputs) == 1 {
		logger.SetOutput(outputs[0])
	}

	defaultLogger = &Logger{Logger: logger}
	return defaultLogger
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		// Initialize with default config if not already initialized
		config := Config{
			Level:  InfoLevel,
			Format: "json",
			Output: "stdout",
		}
		return Initialize(config)
	}
	return defaultLogger
}

// WithCorrelationID adds correlation ID to logger context
func (l *Logger) WithCorrelationID(correlationID string) *logrus.Entry {
	return l.WithField("correlation_id", correlationID)
}

// WithUserID adds user ID to logger context
func (l *Logger) WithUserID(userID uint) *logrus.Entry {
	return l.WithField("user_id", userID)
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField("request_id", requestID)
}

// WithOperation adds operation name to logger context
func (l *Logger) WithOperation(operation string) *logrus.Entry {
	return l.WithField("operation", operation)
}

// WithService adds service name to logger context
func (l *Logger) WithService(service string) *logrus.Entry {
	return l.WithField("service", service)
}

// WithError adds error to logger context
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithFields adds multiple fields to logger context
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithContext extracts logging context from context.Context
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.Logger.WithFields(logrus.Fields{})

	if correlationID := ctx.Value(CorrelationIDKey); correlationID != nil {
		entry = entry.WithField("correlation_id", correlationID)
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	if operation := ctx.Value(OperationKey); operation != nil {
		entry = entry.WithField("operation", operation)
	}

	if service := ctx.Value(ServiceKey); service != nil {
		entry = entry.WithField("service", service)
	}

	return entry
}

// Helper functions for context management

// NewCorrelationID generates a new correlation ID
func NewCorrelationID() string {
	return uuid.New().String()
}

// AddCorrelationIDToContext adds correlation ID to context
func AddCorrelationIDToContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// AddUserIDToContext adds user ID to context
func AddUserIDToContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// AddRequestIDToContext adds request ID to context
func AddRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// AddOperationToContext adds operation to context
func AddOperationToContext(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// AddServiceToContext adds service name to context
func AddServiceToContext(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, ServiceKey, service)
}

// GetCorrelationIDFromContext retrieves correlation ID from context
func GetCorrelationIDFromContext(ctx context.Context) (string, bool) {
	correlationID, ok := ctx.Value(CorrelationIDKey).(string)
	return correlationID, ok
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// Package-level helper functions using the default logger

// Debug logs debug message
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info logs info message
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn logs warning message
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error logs error message
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Fatal logs fatal message and exits
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// WithCorrelationID creates an entry with correlation ID
func WithCorrelationID(correlationID string) *logrus.Entry {
	return GetLogger().WithCorrelationID(correlationID)
}

// WithUserID creates an entry with user ID
func WithUserID(userID uint) *logrus.Entry {
	return GetLogger().WithUserID(userID)
}

// WithOperation creates an entry with operation
func WithOperation(operation string) *logrus.Entry {
	return GetLogger().WithOperation(operation)
}

// WithService creates an entry with service name
func WithService(service string) *logrus.Entry {
	return GetLogger().WithService(service)
}

// WithContext creates an entry with context
func WithContext(ctx context.Context) *logrus.Entry {
	return GetLogger().WithContext(ctx)
}

// WithFields creates an entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithError creates an entry with error
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}
