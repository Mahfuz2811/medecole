package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// WrappedResponse represents the standard wrapped API response
type WrappedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// SuccessResponseWithMessage sends a successful JSON response with custom message
func SuccessResponseWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, WrappedResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorNotFound sends a not found error response
func ErrorNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: message,
		Code:  "NOT_FOUND",
	})
}

// ErrorBadRequest sends a bad request error response
func ErrorBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: message,
		Code:  "BAD_REQUEST",
	})
}

// ErrorInternalServer sends an internal server error response
func ErrorInternalServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: message,
		Code:  "INTERNAL_ERROR",
	})
}

// ErrorValidation sends a validation error response
func ErrorValidation(c *gin.Context, message string, details string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   message,
		Code:    "VALIDATION_ERROR",
		Details: details,
	})
}

// ErrorUnauthorized sends an unauthorized error response
func ErrorUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: message,
		Code:  "UNAUTHORIZED",
	})
}
