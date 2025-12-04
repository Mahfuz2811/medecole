package handlers

import (
	"errors"
	"github.com/Mahfuz2811/medecole/backend/internal/cache"
	"github.com/Mahfuz2811/medecole/backend/internal/dto"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"github.com/Mahfuz2811/medecole/backend/internal/response"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PackageHandler handles package-related HTTP requests
type PackageHandler struct {
	packageService service.PackageService
}

// NewPackageHandler creates a new package handler
func NewPackageHandler(packageService service.PackageService) *PackageHandler {
	return &PackageHandler{
		packageService: packageService,
	}
}

// GetPackages handles GET /api/packages - List all active packages
func (h *PackageHandler) GetPackages(c *gin.Context) {
	// No query parameters needed - get all active packages ordered by sort_order
	var req dto.PackageListRequest

	// Create context from Gin context
	ctx := c.Request.Context()

	updatedCtx, packages, err := h.packageService.GetPackages(ctx, req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch packages")
		return
	}

	// Add cache metadata headers if available
	if metadata := cache.GetCacheMetadata(updatedCtx); metadata != nil {
		c.Header("X-Cache-Status", metadata.Status)
		c.Header("X-Cache-Source", metadata.Source)
		c.Header("X-Cache-TTL", strconv.FormatInt(metadata.TTL, 10))
	}

	response.SuccessResponse(c, packages)
}

// GetPackageBySlug handles GET /api/packages/:slug - Get specific package with exam schedule
func (h *PackageHandler) GetPackageBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// Create context from Gin context
	ctx := c.Request.Context()

	// For package details, include exam schedule data
	updatedCtx, pkg, err := h.packageService.GetPackageBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrPackageNotFound) {
			response.ErrorNotFound(c, "Package not found")
			return
		}
		response.ErrorInternalServer(c, "Failed to fetch package")
		return
	}

	// Add cache metadata headers if available
	if metadata := cache.GetCacheMetadata(updatedCtx); metadata != nil {
		c.Header("X-Cache-Status", metadata.Status)
		c.Header("X-Cache-Source", metadata.Source)
		c.Header("X-Cache-TTL", strconv.FormatInt(metadata.TTL, 10))
	}

	response.SuccessResponse(c, pkg)
}
