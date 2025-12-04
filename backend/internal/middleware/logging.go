package middleware

import (
	"quizora-backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware creates a logging middleware that adds structured logging to requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Create correlation ID for the request
		correlationID := logger.NewCorrelationID()

		// Add correlation ID to context
		param.Request = param.Request.WithContext(
			logger.AddCorrelationIDToContext(param.Request.Context(), correlationID),
		)

		// Extract user ID if available (from auth middleware)
		var userID interface{}
		if param.Keys != nil {
			userID = param.Keys["userID"]
		}

		// Log the request with structured fields
		fields := logrus.Fields{
			"correlation_id":   correlationID,
			"method":           param.Method,
			"path":             param.Path,
			"status_code":      param.StatusCode,
			"latency":          param.Latency,
			"response_time_ms": float64(param.Latency.Nanoseconds()) / 1000000, // Convert to milliseconds
			"client_ip":        param.ClientIP,
			"user_agent":       param.Request.UserAgent(),
			"request_size":     param.Request.ContentLength,
			"response_size":    param.BodySize,
		}

		if userID != nil {
			fields["user_id"] = userID
		}

		// Add error if present
		if param.ErrorMessage != "" {
			fields["error"] = param.ErrorMessage
		}

		// Determine log level based on status code
		logEntry := logger.WithFields(fields)

		switch {
		case param.StatusCode >= 500:
			logEntry.Error("HTTP Request completed with server error")
		case param.StatusCode >= 400:
			logEntry.Warn("HTTP Request completed with client error")
		case param.StatusCode >= 300:
			logEntry.Info("HTTP Request completed with redirect")
		default:
			logEntry.Info("HTTP Request completed successfully")
		}

		// Return empty string as we handle logging ourselves
		return ""
	})
}

// RequestTracingMiddleware adds request tracing context
func RequestTracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate correlation ID if not present
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = logger.NewCorrelationID()
		}

		// Generate request ID
		requestID := logger.NewCorrelationID()

		// Add IDs to context
		ctx := c.Request.Context()
		ctx = logger.AddCorrelationIDToContext(ctx, correlationID)
		ctx = logger.AddRequestIDToContext(ctx, requestID)

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Add headers for response
		c.Header("X-Correlation-ID", correlationID)
		c.Header("X-Request-ID", requestID)

		// Store in gin context for easy access
		c.Set("correlation_id", correlationID)
		c.Set("request_id", requestID)

		// Log request start
		logger.WithFields(logrus.Fields{
			"correlation_id": correlationID,
			"request_id":     requestID,
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"query":          c.Request.URL.RawQuery,
			"client_ip":      c.ClientIP(),
			"user_agent":     c.Request.UserAgent(),
			"content_length": c.Request.ContentLength,
		}).Info("HTTP Request started")

		c.Next()
	}
}

// ErrorLoggingMiddleware logs detailed error information
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			correlationID, _ := c.Get("correlation_id")
			requestID, _ := c.Get("request_id")
			userID, _ := c.Get("userID")

			for _, ginErr := range c.Errors {
				fields := logrus.Fields{
					"correlation_id": correlationID,
					"request_id":     requestID,
					"method":         c.Request.Method,
					"path":           c.Request.URL.Path,
					"error_type":     ginErr.Type,
					"error_message":  ginErr.Error(),
				}

				if userID != nil {
					fields["user_id"] = userID
				}

				logger.WithFields(fields).Error("Request processing error")
			}
		}
	}
}

// PanicRecoveryMiddleware logs panic information with context
func PanicRecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		correlationID, _ := c.Get("correlation_id")
		requestID, _ := c.Get("request_id")
		userID, _ := c.Get("userID")

		fields := logrus.Fields{
			"correlation_id": correlationID,
			"request_id":     requestID,
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"panic":          recovered,
		}

		if userID != nil {
			fields["user_id"] = userID
		}

		logger.WithFields(fields).Error("Panic recovered during request processing")

		c.AbortWithStatus(500)
	})
}
