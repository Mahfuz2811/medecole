package handlers

import (
	"errors"
	"quizora-backend/internal/dto"
	"quizora-backend/internal/logger"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/response"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ExamHandler handles exam-related HTTP requests
type ExamHandler struct {
	examService service.ExamService
}

// NewExamHandler creates a new exam handler
func NewExamHandler(examService service.ExamService) *ExamHandler {
	return &ExamHandler{
		examService: examService,
	}
}

// GetPackageExams handles GET /api/packages/:slug/exams - List all exams for a specific package
func (h *ExamHandler) GetPackageExams(c *gin.Context) {
	packageSlug := c.Param("slug")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Get all package exams (no filtering)
	exams, err := h.examService.GetPackageExamsBySlug(packageSlug, userIDUint)
	if err != nil {
		if errors.Is(err, repository.ErrPackageNotFound) {
			response.ErrorNotFound(c, "Package not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to fetch package exams")
		return
	}

	response.SuccessResponse(c, exams)
}

// GetExamMeta handles GET /api/exams/meta/:slug - Get exam metadata only (no user data or questions)
func (h *ExamHandler) GetExamMeta(c *gin.Context) {
	examSlug := c.Param("slug")

	// Get exam metadata (no user authentication required for metadata)
	examMeta, err := h.examService.GetExamMetaBySlug(examSlug)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			response.ErrorNotFound(c, "Exam not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to fetch exam metadata")
		return
	}

	response.SuccessResponse(c, examMeta)
}

// StartExam handles POST /api/exams/:slug/start - Start a new exam session
func (h *ExamHandler) StartExam(c *gin.Context) {
	examSlug := c.Param("slug")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Parse request body
	var req dto.StartExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request format")
		return
	}

	// Start exam session
	sessionResponse, err := h.examService.StartExam(req.PackageSlug, examSlug, userIDUint, req.DeviceInfo)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			response.ErrorNotFound(c, "Exam not found")
			return
		}
		if errors.Is(err, repository.ErrNotEnrolledInPackage) {
			response.ErrorBadRequest(c, "You must be enrolled in this package to take the exam")
			return
		}
		if errors.Is(err, repository.ErrExamAlreadySubmitted) {
			response.ErrorBadRequest(c, "You have already completed this exam. Multiple attempts are not allowed.")
			return
		}
		response.ErrorInternalServer(c, "Failed to start exam session")
		return
	}

	response.SuccessResponse(c, sessionResponse)
}

// GetSession handles GET /api/exams/session/:sessionId - Get exam session data
func (h *ExamHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Initialize logger with operation context
	log := logger.WithOperation("GetSession").WithFields(logrus.Fields{
		"session_id": sessionID,
	})

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Add user ID to logger context
	log = log.WithField("user_id", userIDUint)

	sessionData, err := h.examService.GetSession(sessionID, userIDUint)
	if err != nil {
		// Log different error types with appropriate levels and context
		if errors.Is(err, repository.ErrAttemptNotFound) {
			log.WithError(err).Warn("Session not found - invalid session ID or user mismatch")
			response.ErrorNotFound(c, "Session not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to fetch session data")
		return
	}

	response.SuccessResponse(c, sessionData)
}

// SyncSession handles PUT /api/exams/session/:sessionId/sync - Sync user answers during exam session
func (h *ExamHandler) SyncSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Parse request body
	var req dto.SyncSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request format")
		return
	}

	// Validate that we have answers to sync
	if len(req.Answers) == 0 {
		response.ErrorBadRequest(c, "No answers provided")
		return
	}

	// Sync session answers
	syncResponse, err := h.examService.SyncSession(sessionID, userIDUint, req.Answers)
	if err != nil {
		if errors.Is(err, repository.ErrAttemptNotFound) {
			response.ErrorNotFound(c, "Session not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to sync session answers")
		return
	}

	response.SuccessResponse(c, syncResponse)
}

// SubmitExam handles POST /api/exams/submit - Submit exam and finalize session
func (h *ExamHandler) SubmitExam(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Parse request body
	var req dto.SubmitExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request format")
		return
	}

	// Submit exam
	submitResponse, err := h.examService.SubmitExam(req.SessionID, userIDUint)
	if err != nil {
		if errors.Is(err, repository.ErrAttemptNotFound) {
			response.ErrorNotFound(c, "Session not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to submit exam")
		return
	}

	response.SuccessResponse(c, submitResponse)
}

// GetExamResults handles GET /api/exams/results/:sessionId - Get exam results by session
func (h *ExamHandler) GetExamResults(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Initialize logger with operation context
	log := logger.WithOperation("GetExamResults").WithFields(logrus.Fields{
		"session_id": sessionID,
	})

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrorUnauthorized(c, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.ErrorUnauthorized(c, "Invalid user ID")
		return
	}

	// Add user ID to logger context
	log = log.WithField("user_id", userIDUint)

	// Get exam results by session ID
	results, err := h.examService.GetExamResultsBySession(sessionID, userIDUint)
	if err != nil {
		// Log different error types with appropriate levels and context
		if errors.Is(err, repository.ErrExamNotFound) {
			log.WithError(err).Warn("Exam session not found - invalid session ID or user mismatch")
			response.ErrorNotFound(c, "Exam session not found")
			return
		}
		log.WithError(err).Error("Failed to retrieve exam results from service")
		response.ErrorInternalServer(c, "Failed to fetch exam results")
		return
	}

	response.SuccessResponseWithMessage(c, results, "Exam results retrieved successfully")
}
