package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"quizora-backend/internal/cache"
	"quizora-backend/internal/dto"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/models"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageRepository is a mock implementation of PackageRepository
type MockPackageRepository struct {
	mock.Mock
}

func (m *MockPackageRepository) GetActivePackages() ([]models.Package, error) {
	args := m.Called()
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockPackageRepository) GetBySlugWithExams(slug string) (*models.Package, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Package), args.Error(1)
}

// MockCache mocks the CacheInterface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCache) Exists(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockCache) GetTTL(key string) (time.Duration, error) {
	args := m.Called(key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCache) Clear() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper functions to create test data
func createTestPackage() models.Package {
	now := time.Now()
	validityDays := 30
	price := 99.99
	description := "Test package description"
	imageURL := "https://example.com/image.jpg"
	imageAlt := "Test package image"

	return models.Package{
		ID:              1,
		Name:            "Test Package",
		Slug:            "test-package",
		Description:     &description,
		PackageType:     models.PackageTypePremium,
		Price:           price,
		ImageURL:        &imageURL,
		ImageAlt:        &imageAlt,
		ValidityType:    models.ValidityTypeRelative,
		ValidityDays:    &validityDays,
		TotalExams:      5,
		IsActive:        true,
		SortOrder:       1,
		CreatedAt:       now,
		UpdatedAt:       now,
		EnrollmentCount: 100,
	}
}

func createTestPackageWithExams() models.Package {
	pkg := createTestPackage()

	// Add package exams
	examTime := time.Now().Add(24 * time.Hour)
	pkg.PackageExams = []models.PackageExam{
		{
			ID:        1,
			PackageID: 1,
			ExamID:    1,
			SortOrder: 1,
			IsActive:  true,
			Exam: models.Exam{
				ID:                    1,
				Title:                 "Test Exam 1",
				Slug:                  "test-exam-1",
				ExamType:              models.ExamTypeDaily,
				TotalQuestions:        10,
				DurationMinutes:       60,
				PassingScore:          60.0,
				MaxAttempts:           3,
				ScheduledStartDate:    &examTime,
				Status:                models.ExamStatusScheduled,
				IsActive:              true,
				AttemptCount:          50,
				CompletedAttemptCount: 45,
			},
		},
	}

	return pkg
}

func createTestPackageList() []models.Package {
	return []models.Package{
		createTestPackage(),
		{
			ID:              2,
			Name:            "Free Package",
			Slug:            "free-package",
			PackageType:     models.PackageTypeFree,
			Price:           0.0,
			ValidityType:    models.ValidityTypeFixed,
			TotalExams:      3,
			IsActive:        true,
			SortOrder:       2,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			EnrollmentCount: 500,
		},
	}
}

func TestPackageService_GetPackages_TTLProgression(t *testing.T) {
	// Setup with real memory cache to test actual TTL behavior
	mockRepo := new(MockPackageRepository)
	memoryCache := cache.NewMemoryCache(50, 1000) // Real memory cache
	defer memoryCache.Close()
	mapper := mapper.NewPackageMapper()
	service := service.NewPackageService(mockRepo, mapper, memoryCache)

	// Mock data
	testPackages := []models.Package{createTestPackage()}
	mockRepo.On("GetActivePackages").Return(testPackages, nil)

	t.Run("TTL decreases over time", func(t *testing.T) {
		// First request - cache miss
		ctx1, result1, err1 := service.GetPackages(context.Background(), dto.PackageListRequest{})
		assert.NoError(t, err1)
		assert.NotNil(t, result1)

		metadata1 := cache.GetCacheMetadata(ctx1)
		assert.NotNil(t, metadata1)
		assert.Equal(t, "MISS", metadata1.Status)
		assert.Equal(t, "database", metadata1.Source)
		assert.Equal(t, int64(0), metadata1.TTL)

		// Wait a small amount to ensure cache is set
		time.Sleep(100 * time.Millisecond)

		// Second request - cache hit with full TTL (approximately 120 seconds)
		ctx2, result2, err2 := service.GetPackages(context.Background(), dto.PackageListRequest{})
		assert.NoError(t, err2)
		assert.NotNil(t, result2)

		metadata2 := cache.GetCacheMetadata(ctx2)
		assert.NotNil(t, metadata2)
		assert.Equal(t, "HIT", metadata2.Status)
		assert.Equal(t, "memory", metadata2.Source)
		// TTL should be close to 120 seconds (within 5 seconds tolerance due to execution time)
		assert.Greater(t, metadata2.TTL, int64(115))
		assert.LessOrEqual(t, metadata2.TTL, int64(120))

		initialTTL := metadata2.TTL

		// Wait 2 seconds
		time.Sleep(2 * time.Second)

		// Third request - cache hit with reduced TTL
		ctx3, result3, err3 := service.GetPackages(context.Background(), dto.PackageListRequest{})
		assert.NoError(t, err3)
		assert.NotNil(t, result3)

		metadata3 := cache.GetCacheMetadata(ctx3)
		assert.NotNil(t, metadata3)
		assert.Equal(t, "HIT", metadata3.Status)
		assert.Equal(t, "memory", metadata3.Source)
		// TTL should be less than initial TTL (allowing for some execution time variance)
		assert.Less(t, metadata3.TTL, initialTTL)
		assert.Greater(t, metadata3.TTL, initialTTL-5) // Should not decrease by more than 5 seconds

		t.Logf("Initial TTL: %d seconds, After 2 seconds: %d seconds", initialTTL, metadata3.TTL)
	})

	mockRepo.AssertExpectations(t)
}

func TestPackageService_GetPackages_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	mockCache := new(MockCache)
	mapper := mapper.NewPackageMapper()
	service := service.NewPackageService(mockRepo, mapper, mockCache)

	// Mock data
	expectedResponse := dto.PackageListResponse{
		Packages: []dto.PackageResponse{},
	}

	// Setup cache expectations (cache hit)
	mockCache.On("Get", "packages:list", mock.AnythingOfType("*dto.PackageListResponse")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate cache hit by populating the destination
		dest := args[1].(*dto.PackageListResponse)
		*dest = expectedResponse
	})
	mockCache.On("GetTTL", "packages:list").Return(90*time.Second, nil) // 90 seconds remaining

	// Execute
	ctx, result, err := service.GetPackages(context.Background(), dto.PackageListRequest{})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, *result)

	// Check cache metadata in context
	metadata := cache.GetCacheMetadata(ctx)
	assert.NotNil(t, metadata)
	assert.Equal(t, "HIT", metadata.Status)
	assert.Equal(t, "memory", metadata.Source)
	assert.Equal(t, int64(90), metadata.TTL) // Should be 90 seconds remaining

	mockCache.AssertExpectations(t)
	// Repository should not be called on cache hit
	mockRepo.AssertNotCalled(t, "GetActivePackages")
}

func TestPackageService_GetPackages(t *testing.T) {
	tests := []struct {
		name           string
		request        dto.PackageListRequest
		mockPackages   []models.Package
		mockError      error
		expectedError  error
		expectedResult *dto.PackageListResponse
	}{
		{
			name:          "Success - Get all active packages",
			request:       dto.PackageListRequest{},
			mockPackages:  createTestPackageList(),
			mockError:     nil,
			expectedError: nil,
			expectedResult: &dto.PackageListResponse{
				Packages: []dto.PackageResponse{}, // Will be populated by mapper
			},
		},
		{
			name:          "Success - Empty result",
			request:       dto.PackageListRequest{},
			mockPackages:  []models.Package{},
			mockError:     nil,
			expectedError: nil,
			expectedResult: &dto.PackageListResponse{
				Packages: []dto.PackageResponse{},
			},
		},
		{
			name:           "Error - Repository error",
			request:        dto.PackageListRequest{},
			mockPackages:   nil,
			mockError:      errors.New("database connection error"),
			expectedError:  errors.New("database connection error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockPackageRepository)
			mockCache := new(MockCache)
			mapper := mapper.NewPackageMapper()
			service := service.NewPackageService(mockRepo, mapper, mockCache)

			// Setup cache expectations (return cache miss to test database fallback)
			mockCache.On("Get", "packages:list", mock.AnythingOfType("*dto.PackageListResponse")).Return(cache.ErrKeyNotFound)
			mockCache.On("Set", "packages:list", mock.AnythingOfType("dto.PackageListResponse"), 2*time.Minute).Return(nil)

			// Setup expectations - no parameters needed for simplified method
			mockRepo.On("GetActivePackages").Return(tt.mockPackages, tt.mockError)

			// Execute
			_, result, err := service.GetPackages(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Packages, len(tt.mockPackages))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPackageService_GetPackageBySlug_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	mockCache := new(MockCache)
	mapper := mapper.NewPackageMapper()
	service := service.NewPackageService(mockRepo, mapper, mockCache)

	// Mock data
	expectedResponse := dto.PackageResponse{
		ID:   1,
		Name: "Test Package",
		Slug: "test-package",
	}

	// Setup cache expectations (cache hit)
	mockCache.On("Get", "package:slug:test-package", mock.AnythingOfType("*dto.PackageResponse")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate cache hit by populating the destination
		dest := args[1].(*dto.PackageResponse)
		*dest = expectedResponse
	})
	mockCache.On("GetTTL", "package:slug:test-package").Return(240*time.Second, nil) // 4 minutes remaining

	// Execute
	ctx, result, err := service.GetPackageBySlug(context.Background(), "test-package")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, *result)

	// Check cache metadata in context
	metadata := cache.GetCacheMetadata(ctx)
	assert.NotNil(t, metadata)
	assert.Equal(t, "HIT", metadata.Status)
	assert.Equal(t, "memory", metadata.Source)
	assert.Equal(t, int64(240), metadata.TTL) // Should be 240 seconds remaining

	mockCache.AssertExpectations(t)
	// Repository should not be called on cache hit
	mockRepo.AssertNotCalled(t, "GetBySlugWithExams")
}

func TestPackageService_GetPackageBySlug_WithExams_CacheHit(t *testing.T) {
	// Setup
	mockRepo := new(MockPackageRepository)
	mockCache := new(MockCache)
	mapper := mapper.NewPackageMapper()
	service := service.NewPackageService(mockRepo, mapper, mockCache)

	// Mock data
	expectedResponse := dto.PackageResponse{
		ID:    1,
		Name:  "Test Package",
		Slug:  "test-package",
		Exams: []dto.PackageExamScheduleResponse{}, // With exams data
	}

	// Setup cache expectations (cache hit)
	mockCache.On("Get", "package:slug:test-package", mock.AnythingOfType("*dto.PackageResponse")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate cache hit by populating the destination
		dest := args[1].(*dto.PackageResponse)
		*dest = expectedResponse
	})
	mockCache.On("GetTTL", "package:slug:test-package").Return(180*time.Second, nil) // 3 minutes remaining

	// Execute
	ctx, result, err := service.GetPackageBySlug(context.Background(), "test-package")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, *result)

	// Check cache metadata in context
	metadata := cache.GetCacheMetadata(ctx)
	assert.NotNil(t, metadata)
	assert.Equal(t, "HIT", metadata.Status)
	assert.Equal(t, "memory", metadata.Source)
	assert.Equal(t, int64(180), metadata.TTL) // Should be 180 seconds remaining

	mockCache.AssertExpectations(t)
	// Repository should not be called on cache hit
	mockRepo.AssertNotCalled(t, "GetBySlugWithExams")
}

func TestPackageService_GetPackageBySlug(t *testing.T) {
	tests := []struct {
		name           string
		slug           string
		mockPackage    *models.Package
		mockError      error
		expectedError  error
		expectedResult bool // whether result should be non-nil
	}{
		{
			name:           "Success - Valid slug",
			slug:           "test-package",
			mockPackage:    func() *models.Package { pkg := createTestPackage(); return &pkg }(),
			mockError:      nil,
			expectedError:  nil,
			expectedResult: true,
		},
		{
			name:           "Error - Package not found",
			slug:           "non-existent-package",
			mockPackage:    nil,
			mockError:      repository.ErrPackageNotFound,
			expectedError:  repository.ErrPackageNotFound,
			expectedResult: false,
		},
		{
			name:           "Error - Database error",
			slug:           "test-package",
			mockPackage:    nil,
			mockError:      errors.New("database connection error"),
			expectedError:  errors.New("database connection error"),
			expectedResult: false,
		},
		{
			name:           "Error - Empty slug",
			slug:           "",
			mockPackage:    nil,
			mockError:      repository.ErrPackageNotFound,
			expectedError:  repository.ErrPackageNotFound,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockPackageRepository)
			mockCache := new(MockCache)
			mapper := mapper.NewPackageMapper()
			service := service.NewPackageService(mockRepo, mapper, mockCache)

			// Setup expectations
			mockCache.On("Get", "package:slug:"+tt.slug, mock.AnythingOfType("*dto.PackageResponse")).Return(cache.ErrKeyNotFound)
			if tt.mockError == nil {
				// Only expect Set call if there's no repository error
				mockCache.On("Set", "package:slug:"+tt.slug, mock.AnythingOfType("dto.PackageResponse"), 5*time.Minute).Return(nil)
			}
			mockRepo.On("GetBySlugWithExams", tt.slug).Return(tt.mockPackage, tt.mockError)

			// Execute
			_, result, err := service.GetPackageBySlug(context.Background(), tt.slug)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedResult {
					assert.NotNil(t, result)
					assert.Equal(t, tt.mockPackage.ID, result.ID)
					assert.Equal(t, tt.mockPackage.Name, result.Name)
					assert.Equal(t, tt.mockPackage.Slug, result.Slug)
					assert.Equal(t, tt.mockPackage.PackageType, result.PackageType)
					assert.Equal(t, tt.mockPackage.Price, result.Price)
					assert.Equal(t, tt.mockPackage.IsActive, result.IsActive)
					// Verify that exams array is present (method now includes exams)
					assert.NotNil(t, result.Exams)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPackageService_GetPackageBySlug_WithExams(t *testing.T) {
	tests := []struct {
		name           string
		slug           string
		mockPackage    *models.Package
		mockError      error
		expectedError  error
		expectedResult bool
	}{
		{
			name:           "Success - Package with exams",
			slug:           "test-package",
			mockPackage:    func() *models.Package { pkg := createTestPackageWithExams(); return &pkg }(),
			mockError:      nil,
			expectedError:  nil,
			expectedResult: true,
		},
		{
			name:           "Success - Package without exams",
			slug:           "empty-package",
			mockPackage:    func() *models.Package { pkg := createTestPackage(); return &pkg }(),
			mockError:      nil,
			expectedError:  nil,
			expectedResult: true,
		},
		{
			name:           "Error - Package not found",
			slug:           "non-existent-package",
			mockPackage:    nil,
			mockError:      repository.ErrPackageNotFound,
			expectedError:  repository.ErrPackageNotFound,
			expectedResult: false,
		},
		{
			name:           "Error - Database error",
			slug:           "test-package",
			mockPackage:    nil,
			mockError:      errors.New("failed to preload exams"),
			expectedError:  errors.New("failed to preload exams"),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockPackageRepository)
			mockCache := new(MockCache)
			mapper := mapper.NewPackageMapper()
			service := service.NewPackageService(mockRepo, mapper, mockCache)

			// Setup expectations
			mockCache.On("Get", "package:slug:"+tt.slug, mock.AnythingOfType("*dto.PackageResponse")).Return(cache.ErrKeyNotFound)
			if tt.mockError == nil {
				// Only expect Set call if there's no repository error
				mockCache.On("Set", "package:slug:"+tt.slug, mock.AnythingOfType("dto.PackageResponse"), 5*time.Minute).Return(nil)
			}
			mockRepo.On("GetBySlugWithExams", tt.slug).Return(tt.mockPackage, tt.mockError)

			// Execute
			_, result, err := service.GetPackageBySlug(context.Background(), tt.slug)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedResult {
					assert.NotNil(t, result)
					assert.Equal(t, tt.mockPackage.ID, result.ID)
					assert.Equal(t, tt.mockPackage.Name, result.Name)
					assert.Equal(t, tt.mockPackage.Slug, result.Slug)

					// Verify exams are included
					if tt.mockPackage.PackageExams != nil {
						assert.Len(t, result.Exams, len(tt.mockPackage.PackageExams))
						if len(result.Exams) > 0 {
							assert.Equal(t, tt.mockPackage.PackageExams[0].Exam.Title, result.Exams[0].Exam.Title)
						}
					} else {
						assert.Empty(t, result.Exams)
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPackageService_Integration(t *testing.T) {
	t.Run("Service methods work together", func(t *testing.T) {
		// Setup
		mockRepo := new(MockPackageRepository)
		mockCache := new(MockCache)
		mapper := mapper.NewPackageMapper()
		service := service.NewPackageService(mockRepo, mapper, mockCache)

		// Test data
		testPackage := createTestPackage()
		testPackageWithExams := createTestPackageWithExams()

		// Setup expectations for multiple calls
		mockCache.On("Get", "packages:list", mock.AnythingOfType("*dto.PackageListResponse")).Return(cache.ErrKeyNotFound)
		mockCache.On("Set", "packages:list", mock.AnythingOfType("dto.PackageListResponse"), 2*time.Minute).Return(nil)

		// Cache expectations for GetPackageBySlug (now returns exams)
		mockCache.On("Get", "package:slug:test-package", mock.AnythingOfType("*dto.PackageResponse")).Return(cache.ErrKeyNotFound)
		mockCache.On("Set", "package:slug:test-package", mock.AnythingOfType("dto.PackageResponse"), 5*time.Minute).Return(nil)

		mockRepo.On("GetActivePackages").Return([]models.Package{testPackage}, nil)
		mockRepo.On("GetBySlugWithExams", "test-package").Return(&testPackageWithExams, nil)

		// Test GetPackages
		_, listResult, err := service.GetPackages(context.Background(), dto.PackageListRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, listResult)
		assert.Len(t, listResult.Packages, 1)

		// Test GetPackageBySlug (now includes exams)
		_, singleResult, err := service.GetPackageBySlug(context.Background(), "test-package")
		assert.NoError(t, err)
		assert.NotNil(t, singleResult)
		assert.NotEmpty(t, singleResult.Exams)
		assert.Equal(t, "Test Exam 1", singleResult.Exams[0].Exam.Title)

		mockRepo.AssertExpectations(t)
	})
}

func TestPackageService_Constructor(t *testing.T) {
	t.Run("NewPackageService creates service correctly", func(t *testing.T) {
		mockRepo := new(MockPackageRepository)
		mockCache := new(MockCache)
		mapper := mapper.NewPackageMapper()

		service := service.NewPackageService(mockRepo, mapper, mockCache)

		assert.NotNil(t, service)
	})

	t.Run("Service interface compliance", func(t *testing.T) {
		mockRepo := new(MockPackageRepository)
		mockCache := new(MockCache)
		mapper := mapper.NewPackageMapper()

		svc := service.NewPackageService(mockRepo, mapper, mockCache)

		// Verify service implements the interface by checking it can be assigned
		var _ service.PackageService = svc
		assert.NotNil(t, svc)
	})
}

func TestPackageService_PackageDetailsCaching_Integration(t *testing.T) {
	// Setup with real memory cache to test actual caching behavior
	mockRepo := new(MockPackageRepository)
	memoryCache := cache.NewMemoryCache(50, 1000) // Real memory cache
	defer memoryCache.Close()
	mapper := mapper.NewPackageMapper()
	service := service.NewPackageService(mockRepo, mapper, memoryCache)

	// Mock data
	testPackage := createTestPackageWithExams()
	mockRepo.On("GetBySlugWithExams", "test-package").Return(&testPackage, nil)

	t.Run("Package details caching with TTL progression", func(t *testing.T) {
		// First request - cache miss
		ctx1, result1, err1 := service.GetPackageBySlug(context.Background(), "test-package")
		assert.NoError(t, err1)
		assert.NotNil(t, result1)

		metadata1 := cache.GetCacheMetadata(ctx1)
		assert.NotNil(t, metadata1)
		assert.Equal(t, "MISS", metadata1.Status)
		assert.Equal(t, "database", metadata1.Source)
		assert.Equal(t, int64(0), metadata1.TTL)

		// Wait a small amount to ensure cache is set
		time.Sleep(100 * time.Millisecond)

		// Second request - cache hit with full TTL (approximately 300 seconds for 5 minutes)
		ctx2, result2, err2 := service.GetPackageBySlug(context.Background(), "test-package")
		assert.NoError(t, err2)
		assert.NotNil(t, result2)

		metadata2 := cache.GetCacheMetadata(ctx2)
		assert.NotNil(t, metadata2)
		assert.Equal(t, "HIT", metadata2.Status)
		assert.Equal(t, "memory", metadata2.Source)
		// TTL should be close to 300 seconds (within 10 seconds tolerance due to execution time)
		assert.Greater(t, metadata2.TTL, int64(290))
		assert.LessOrEqual(t, metadata2.TTL, int64(300))

		initialTTL := metadata2.TTL

		// Wait 2 seconds
		time.Sleep(2 * time.Second)

		// Third request - cache hit with reduced TTL
		ctx3, result3, err3 := service.GetPackageBySlug(context.Background(), "test-package")
		assert.NoError(t, err3)
		assert.NotNil(t, result3)

		metadata3 := cache.GetCacheMetadata(ctx3)
		assert.NotNil(t, metadata3)
		assert.Equal(t, "HIT", metadata3.Status)
		assert.Equal(t, "memory", metadata3.Source)
		// TTL should be less than initial TTL
		assert.Less(t, metadata3.TTL, initialTTL)
		assert.Greater(t, metadata3.TTL, initialTTL-5) // Should not decrease by more than 5 seconds

		t.Logf("Package details cache - Initial TTL: %d seconds, After 2 seconds: %d seconds", initialTTL, metadata3.TTL)

		// Verify the response content is the same for all requests
		assert.Equal(t, result1.ID, result2.ID)
		assert.Equal(t, result1.ID, result3.ID)
		assert.Equal(t, result1.Name, result2.Name)
		assert.Equal(t, result1.Name, result3.Name)
		assert.NotEmpty(t, result3.Exams) // Should include exams
	})

	mockRepo.AssertExpectations(t)
}

func TestPackageService_EdgeCases(t *testing.T) {
	t.Run("Empty packages result", func(t *testing.T) {
		mockRepo := new(MockPackageRepository)
		mockCache := new(MockCache)
		mapper := mapper.NewPackageMapper()
		service := service.NewPackageService(mockRepo, mapper, mockCache)

		request := dto.PackageListRequest{}

		mockCache.On("Get", "packages:list", mock.AnythingOfType("*dto.PackageListResponse")).Return(cache.ErrKeyNotFound)
		mockCache.On("Set", "packages:list", mock.AnythingOfType("dto.PackageListResponse"), 2*time.Minute).Return(nil)
		mockRepo.On("GetActivePackages").Return([]models.Package{}, nil)

		_, result, err := service.GetPackages(context.Background(), request)
		assert.NoError(t, err)
		assert.Empty(t, result.Packages)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Special characters in slug", func(t *testing.T) {
		mockRepo := new(MockPackageRepository)
		mockCache := new(MockCache)
		mapper := mapper.NewPackageMapper()
		service := service.NewPackageService(mockRepo, mapper, mockCache)

		specialSlug := "test-package-with-special-chars-123"

		mockCache.On("Get", "package:slug:"+specialSlug, mock.AnythingOfType("*dto.PackageResponse")).Return(cache.ErrKeyNotFound)
		mockRepo.On("GetBySlugWithExams", specialSlug).Return(nil, repository.ErrPackageNotFound)

		_, result, err := service.GetPackageBySlug(context.Background(), specialSlug)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}
