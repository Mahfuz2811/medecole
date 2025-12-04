package handlers

import (
	"net/http"
	"strconv"

	"quizora-backend/internal/dto"
	"quizora-backend/internal/errors"
	"quizora-backend/internal/logger"
	"quizora-backend/internal/models"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// EnrollmentHandler handles enrollment-related HTTP requests
type EnrollmentHandler struct {
	enrollmentService service.EnrollmentService
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(enrollmentService service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

func (h *EnrollmentHandler) EnrollInPackage(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"handler":   "EnrollmentHandler",
		"operation": "EnrollInPackage",
	})

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid user ID",
		})
		return
	}

	// Parse request body
	var req dto.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Failed to parse enrollment request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	if req.CouponCode != nil {
		log = log.WithField("coupon_code", *req.CouponCode)
	}

	// Additional validation
	if req.PackageID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Package ID is required",
		})
		return
	}

	// Validate coupon code format if provided
	if req.CouponCode != nil && len(*req.CouponCode) > 50 {
		log.Warn("Enrollment request with invalid coupon code length")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Coupon code too long",
		})
		return
	}

	// Process enrollment
	response, err := h.enrollmentService.EnrollInPackage(ctx, uid, req)
	if err != nil {
		log.WithError(err).Error("Enrollment failed")

		// Handle specific business logic errors using type assertions
		switch {
		case errors.IsPackageNotActiveError(err):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Bad Request",
				Message: "Package is not available for enrollment",
			})
			return
		case errors.IsActiveEnrollmentExistsError(err):
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "Conflict",
				Message: "You are already enrolled in this package",
			})
			return
		case errors.IsCouponValidationError(err):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid or expired coupon",
			})
			return
		case errors.IsPackageNotFoundError(err):
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Package not found",
			})
			return
		case errors.IsEnrollmentCreationError(err):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to create enrollment",
			})
			return
		case errors.IsCouponProcessingError(err):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to process coupon",
			})
			return
		case errors.IsPackageStatsUpdateError(err), errors.IsTransactionCommitError(err), errors.IsEnrollmentFetchError(err):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to complete enrollment process",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to process enrollment",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, response)
}

func (h *EnrollmentHandler) CheckEnrollmentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"handler":   "EnrollmentHandler",
		"operation": "CheckEnrollmentStatus",
	})

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		log.Error("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid user ID",
		})
		return
	}

	log = log.WithField("user_id", uid)

	// Parse package ID from query parameter
	packageIDStr := c.Query("package_id")
	if packageIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Package ID is required",
		})
		return
	}

	packageID, err := strconv.ParseUint(packageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid package ID",
		})
		return
	}

	log = log.WithField("package_id", packageID)

	// Check enrollment status
	response, err := h.enrollmentService.CheckEnrollmentStatus(ctx, uid, uint(packageID))
	if err != nil {
		log.WithError(err).Error("Failed to check enrollment status")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to check enrollment status",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EnrollmentHandler) ValidateCoupon(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"handler":   "EnrollmentHandler",
		"operation": "ValidateCoupon",
	})

	log.Debug("Processing coupon validation request")

	// Check if user is authenticated (optional for coupon validation, but good for security)
	_, exists := c.Get("userID")
	if !exists {
		log.Warn("Coupon validation attempt without authentication")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Parse request body
	var req dto.CouponValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Failed to parse coupon validation request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	log = log.WithFields(logrus.Fields{
		"coupon_code": req.CouponCode,
		"package_id":  req.PackageID,
	})

	log.Info("Processing coupon validation request")

	// Additional validation
	if req.PackageID == 0 {
		log.Warn("Coupon validation request missing package ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Package ID is required",
		})
		return
	}

	if req.CouponCode == "" {
		log.Warn("Coupon validation request missing coupon code")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Coupon code is required",
		})
		return
	}

	if len(req.CouponCode) > 50 {
		log.Warn("Coupon validation request with invalid coupon code length")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Coupon code too long",
		})
		return
	}

	// Validate coupon
	response, err := h.enrollmentService.ValidateCoupon(ctx, req)
	if err != nil {
		log.WithError(err).Error("Coupon validation failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to validate coupon",
		})
		return
	}

	log.WithField("coupon_valid", response.Valid).Info("Coupon validation completed")
	c.JSON(http.StatusOK, response)
}
