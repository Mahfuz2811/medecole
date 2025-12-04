package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Mahfuz2811/medecole/backend/internal/cache"
	"github.com/Mahfuz2811/medecole/backend/internal/logger"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ExamRepository handles database operations for exams
type ExamRepository interface {
	GetExamsByPackageSlug(packageSlug string, userID uint) ([]ExamWithUserData, error)
	GetPackageWithExamsBySlug(packageSlug string, userID uint) (*PackageWithExamsData, error)
	GetExamBySlug(examSlug string) (*models.Exam, error)
	GetActiveAttemptByUserAndExam(userID uint, examID uint) (*models.UserExamAttempt, error)
	GetUserAttemptsByExam(userID uint, examID uint) ([]models.UserExamAttempt, error)
	CreateExamAttempt(userID uint, examID uint, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error)
	GetActiveSessionByID(sessionID string) (*SessionWithExamData, error)
	GetCompletedSessionByID(sessionID string) (*SessionWithExamData, error)
	GetSessionAnswers(sessionID string) (map[uint]string, error)
	SyncSessionAnswers(sessionID string, answers map[uint]string) error
	CompleteExamAttempt(attemptID uint, score float64, passed bool) error
	CompleteExamAttemptWithAnswers(attemptID uint, score float64, passed bool, answersData string, correctAnswers int) error
	GetAttemptBySessionAndUser(sessionID string, userID uint) (*models.UserExamAttempt, error)
	MarkExpiredSessionsAsAbandoned(currentTime time.Time, gracePeriodSeconds int) (int64, error)

	// Optimized methods for Phase 1 & 2
	GetUserAttemptForExam(userID uint, examID uint) (*models.UserExamAttempt, error)
	GetUserAttemptForExamInPackage(userID uint, examID uint, packageID uint) (*models.UserExamAttempt, error)
	CreateExamAttemptWithExam(userID uint, exam *models.Exam, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error)
	GetPackageIDForExam(examID uint) (uint, error)
}

// ExamWithUserData represents exam data combined with user attempt information
type ExamWithUserData struct {
	models.Exam
	HasAttempted   bool   `json:"has_attempted"`
	SortOrder      int    `json:"sort_order"`      // From package_exams table
	ComputedStatus string `json:"computed_status"` // Computed based on scheduling
	// User attempt data
	UserAttemptID          *uint      `json:"user_attempt_id,omitempty"`
	UserAttemptStatus      *string    `json:"user_attempt_status,omitempty"` // STARTED, COMPLETED, AUTO_SUBMITTED, ABANDONED
	UserAttemptStartedAt   *time.Time `json:"user_attempt_started_at,omitempty"`
	UserAttemptCompletedAt *time.Time `json:"user_attempt_completed_at,omitempty"`
	AttemptScore           *float64   `json:"attempt_score,omitempty"`
	AttemptCorrectAnswers  *int       `json:"attempt_correct_answers,omitempty"`
	AttemptPassed          *bool      `json:"attempt_passed,omitempty"`
	ActualTimeSpent        *int       `json:"actual_time_spent,omitempty"`
	SessionID              *string    `json:"session_id,omitempty"`
}

// PackageWithExamsData represents package data with its exams and user data
type PackageWithExamsData struct {
	Package models.Package     `json:"package"`
	Exams   []ExamWithUserData `json:"exams"`
}

// SessionWithExamData represents session attempt data with complete exam information
type SessionWithExamData struct {
	Attempt models.UserExamAttempt `json:"attempt"`
	Exam    models.Exam            `json:"exam"`
}

// examRepository implements ExamRepository
type examRepository struct {
	db    *gorm.DB
	cache cache.CacheInterface
}

// NewExamRepository creates a new exam repository
func NewExamRepository(db *gorm.DB, cache cache.CacheInterface) ExamRepository {
	return &examRepository{
		db:    db,
		cache: cache,
	}
}

// GetExamsByPackageSlug retrieves all exams for a specific package by package slug
func (r *examRepository) GetExamsByPackageSlug(packageSlug string, userID uint) ([]ExamWithUserData, error) {
	baseQuery := `
		SELECT 
			e.id,
			e.title,
			e.slug,
			e.description,
			e.exam_type,
			e.total_questions,
			e.duration_minutes,
			e.total_marks,
			e.passing_score,
			e.scheduled_start_date,
			e.scheduled_end_date,
			e.attempt_count,
			e.average_score,
			e.pass_rate,
			CASE 
				WHEN NOW() < e.scheduled_start_date THEN 'UPCOMING'
				WHEN NOW() BETWEEN e.scheduled_start_date AND COALESCE(e.scheduled_end_date, DATE_ADD(e.scheduled_start_date, INTERVAL e.duration_minutes MINUTE)) THEN 'LIVE'
				WHEN e.scheduled_end_date IS NOT NULL AND NOW() > e.scheduled_end_date THEN 'COMPLETED'
				WHEN e.scheduled_start_date IS NOT NULL AND e.scheduled_end_date IS NULL AND NOW() > DATE_ADD(e.scheduled_start_date, INTERVAL e.duration_minutes MINUTE) THEN 'COMPLETED'
				ELSE 'AVAILABLE'
			END as computed_status,
			pe.sort_order,
			CASE WHEN uea.user_id IS NOT NULL THEN true ELSE false END as has_attempted,
			-- User attempt details
			uea.id as user_attempt_id,
			uea.status as user_attempt_status,
			uea.started_at as user_attempt_started_at,
			uea.completed_at as user_attempt_completed_at,
			uea.score as attempt_score,
			uea.correct_answers as attempt_correct_answers,
			uea.is_passed as attempt_passed,
			uea.actual_time_spent as actual_time_spent,
			uea.session_id as session_id
		FROM exams e
		INNER JOIN package_exams pe ON e.id = pe.exam_id
		INNER JOIN packages p ON pe.package_id = p.id
		LEFT JOIN user_exam_attempts uea ON e.id = uea.exam_id AND uea.user_id = ? AND uea.package_id = p.id
		WHERE p.slug = ?
			AND e.is_active = true 
			AND pe.is_active = true
			AND p.is_active = true
		ORDER BY pe.sort_order ASC, e.created_at DESC`

	// Execute query
	var examsWithUserData []ExamWithUserData
	if err := r.db.Raw(baseQuery, userID, packageSlug).Scan(&examsWithUserData).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exams: %w", err)
	}

	return examsWithUserData, nil
}

// GetPackageWithExamsBySlug retrieves package data along with all exams for a specific package by package slug
func (r *examRepository) GetPackageWithExamsBySlug(packageSlug string, userID uint) (*PackageWithExamsData, error) {
	// First get the package data
	var pkg models.Package
	if err := r.db.Where("slug = ? AND is_active = true", packageSlug).First(&pkg).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch package: %w", err)
	}

	// Then get the exams for this package
	exams, err := r.GetExamsByPackageSlug(packageSlug, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exams: %w", err)
	}

	return &PackageWithExamsData{
		Package: pkg,
		Exams:   exams,
	}, nil
}

// GetExamBySlug retrieves an exam by its slug
func (r *examRepository) GetExamBySlug(examSlug string) (*models.Exam, error) {
	var exam models.Exam
	if err := r.db.Where("slug = ? AND is_active = true", examSlug).First(&exam).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrExamNotFound
		}
		return nil, fmt.Errorf("failed to fetch exam: %w", err)
	}
	return &exam, nil
}

// GetActiveAttemptByUserAndExam retrieves an active attempt for a user and exam
func (r *examRepository) GetActiveAttemptByUserAndExam(userID uint, examID uint) (*models.UserExamAttempt, error) {
	var attempt models.UserExamAttempt
	err := r.db.Where("user_id = ? AND exam_id = ? AND status = ?", userID, examID, models.AttemptStatusStarted).First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAttemptNotFound
		}
		return nil, fmt.Errorf("failed to fetch active attempt: %w", err)
	}
	return &attempt, nil
}

// GetUserAttemptsByExam retrieves all attempts for a user and exam
func (r *examRepository) GetUserAttemptsByExam(userID uint, examID uint) ([]models.UserExamAttempt, error) {
	var attempts []models.UserExamAttempt
	err := r.db.Where("user_id = ? AND exam_id = ?", userID, examID).
		Order("created_at DESC").
		Find(&attempts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user attempts: %w", err)
	}
	return attempts, nil
}

// CreateExamAttempt creates a new exam attempt with session
func (r *examRepository) CreateExamAttempt(userID uint, examID uint, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error) {
	// Get exam details for snapshot
	var exam models.Exam
	if err := r.db.First(&exam, examID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exam for attempt: %w", err)
	}

	// Generate session ID
	sessionID := r.generateSessionID()
	now := time.Now()

	// Create attempt record
	attempt := &models.UserExamAttempt{
		UserID:           userID,
		ExamID:           examID,
		PackageID:        packageID,
		AttemptNumber:    1,
		Status:           models.AttemptStatusStarted,
		StartedAt:        now,
		SessionID:        &sessionID,
		LastActivityAt:   &now,
		TimeLimitSeconds: exam.DurationMinutes * 60,
		TotalQuestions:   exam.TotalQuestions,
		PassingScore:     exam.PassingScore,
		AnswersData:      "{}",
		IsScored:         false,
	}

	if err := r.db.Create(attempt).Error; err != nil {
		return nil, fmt.Errorf("failed to create exam attempt: %w", err)
	}

	return attempt, nil
}

// generateSessionID creates a unique session identifier
func (r *examRepository) generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetActiveSessionByID retrieves an active exam session (STARTED status)
// Used during exam execution: GetSession, SyncSession, SubmitExam
func (r *examRepository) GetActiveSessionByID(sessionID string) (*SessionWithExamData, error) {
	// Initialize logger with repository context
	log := logger.WithService("ExamRepository").WithFields(logrus.Fields{
		"operation":  "GetActiveSessionByID",
		"session_id": sessionID,
	})

	var attempt models.UserExamAttempt

	// Get the attempt by session ID (only STARTED status for active sessions)
	err := r.db.Where("session_id = ? AND status = ?", sessionID, models.AttemptStatusStarted).First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAttemptNotFound
		}
		log.WithError(err).Error("Database error while fetching active session attempt")
		return nil, fmt.Errorf("failed to fetch active session attempt: %w", err)
	}

	// Get the associated exam
	var exam models.Exam
	err = r.db.First(&exam, attempt.ExamID).Error
	if err != nil {
		log.WithError(err).WithField("exam_id", attempt.ExamID).Error("Failed to retrieve exam data for active session")
		return nil, fmt.Errorf("failed to fetch exam for active session: %w", err)
	}

	sessionData := &SessionWithExamData{
		Attempt: attempt,
		Exam:    exam,
	}

	return sessionData, nil
}

// GetCompletedSessionByID retrieves a completed exam session (final states only)
// Used for viewing exam results: GetResultsBySession
func (r *examRepository) GetCompletedSessionByID(sessionID string) (*SessionWithExamData, error) {
	// Initialize logger with repository context
	log := logger.WithService("ExamRepository").WithFields(logrus.Fields{
		"operation":  "GetCompletedSessionByID",
		"session_id": sessionID,
	})

	var attempt models.UserExamAttempt

	// Get the attempt by session ID (only allow final states for result viewing)
	err := r.db.Where("session_id = ? AND status IN (?)", sessionID, []models.AttemptStatus{
		models.AttemptStatusCompleted,
		models.AttemptStatusAutoSubmitted,
		models.AttemptStatusAbandoned,
	}).First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAttemptNotFound
		}
		log.WithError(err).Error("Database error while fetching completed session attempt")
		return nil, fmt.Errorf("failed to fetch completed session attempt: %w", err)
	}

	// Get the associated exam
	var exam models.Exam
	err = r.db.First(&exam, attempt.ExamID).Error
	if err != nil {
		log.WithError(err).WithField("exam_id", attempt.ExamID).Error("Failed to retrieve exam data for completed session")
		return nil, fmt.Errorf("failed to fetch exam for completed session: %w", err)
	}

	sessionData := &SessionWithExamData{
		Attempt: attempt,
		Exam:    exam,
	}

	return sessionData, nil
}

// SyncSessionAnswers stores user answers in cache during exam session
func (r *examRepository) SyncSessionAnswers(sessionID string, answers map[uint]string) error {
	// Create cache key for session answers
	cacheKey := fmt.Sprintf("exam_session:%s:answers", sessionID)

	// Store answers in cache with 2 hour TTL (should cover most exam durations)
	err := r.cache.Set(cacheKey, answers, 2*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to sync session answers to cache: %w", err)
	}

	return nil
}

// GetSessionAnswers retrieves user answers from cache for a given session
func (r *examRepository) GetSessionAnswers(sessionID string) (map[uint]string, error) {
	// Create cache key for session answers
	cacheKey := fmt.Sprintf("exam_session:%s:answers", sessionID)

	// Retrieve answers from cache
	var answers map[uint]string
	err := r.cache.Get(cacheKey, &answers)
	if err != nil {
		// If not found in cache, return empty map (not an error - user might not have answered yet)
		return make(map[uint]string), nil
	}

	return answers, nil
}

// CompleteExamAttempt marks an exam attempt as completed and updates the score
func (r *examRepository) CompleteExamAttempt(attemptID uint, score float64, passed bool) error {
	now := time.Now()

	// Get the attempt to calculate time spent
	var attempt models.UserExamAttempt
	err := r.db.First(&attempt, attemptID).Error
	if err != nil {
		return fmt.Errorf("failed to find attempt: %w", err)
	}

	// Calculate actual time spent
	timeSpent := int(now.Sub(attempt.StartedAt).Seconds())
	if timeSpent > attempt.TimeLimitSeconds {
		timeSpent = attempt.TimeLimitSeconds // Cap at time limit for auto-submitted exams
	}

	// Update attempt with completion data
	updates := map[string]interface{}{
		"status":            models.AttemptStatusCompleted,
		"completed_at":      now,
		"actual_time_spent": timeSpent,
		"score":             score,
		"is_scored":         true,
		"session_id":        nil, // Clear session ID as exam is completed
		"last_activity_at":  now,
	}

	err = r.db.Model(&attempt).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to complete exam attempt: %w", err)
	}

	return nil
}

// CompleteExamAttemptWithAnswers marks an exam attempt as completed and updates score with detailed answers
func (r *examRepository) CompleteExamAttemptWithAnswers(attemptID uint, score float64, passed bool, answersData string, correctAnswers int) error {
	now := time.Now()

	// Get the attempt to calculate time spent
	var attempt models.UserExamAttempt
	err := r.db.First(&attempt, attemptID).Error
	if err != nil {
		return fmt.Errorf("failed to find attempt: %w", err)
	}

	// Calculate actual time spent
	timeSpent := int(now.Sub(attempt.StartedAt).Seconds())
	if timeSpent > attempt.TimeLimitSeconds {
		timeSpent = attempt.TimeLimitSeconds // Cap at time limit for auto-submitted exams
	}

	// Update attempt with completion data including answers
	updates := map[string]interface{}{
		"status":            models.AttemptStatusCompleted,
		"completed_at":      now,
		"actual_time_spent": timeSpent,
		"score":             score,
		"correct_answers":   correctAnswers,
		"is_passed":         passed,
		"is_scored":         true,
		"answers_data":      answersData,
		// Keep session_id for tracking purposes - don't clear it
		"last_activity_at": now,
	}

	err = r.db.Model(&attempt).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to complete exam attempt with answers: %w", err)
	}

	return nil
}

// GetAttemptBySessionAndUser gets a user exam attempt by session ID and user ID
func (r *examRepository) GetAttemptBySessionAndUser(sessionID string, userID uint) (*models.UserExamAttempt, error) {
	var attempt models.UserExamAttempt
	err := r.db.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&attempt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExamNotFound
		}
		return nil, fmt.Errorf("failed to get attempt by session: %w", err)
	}
	return &attempt, nil
}

// GetUserAttemptForExam retrieves any attempt (active or completed) for a user and exam
// Optimized single query replacement for GetUserAttemptsByExam + GetActiveAttemptByUserAndExam
func (r *examRepository) GetUserAttemptForExam(userID uint, examID uint) (*models.UserExamAttempt, error) {
	var attempt models.UserExamAttempt
	err := r.db.Where("user_id = ? AND exam_id = ?", userID, examID).
		Order("created_at DESC").
		First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No attempt found, not an error
		}
		return nil, fmt.Errorf("failed to fetch user attempt: %w", err)
	}
	return &attempt, nil
}

// GetUserAttemptForExamInPackage retrieves any attempt (active or completed) for a user, exam, and package
// This ensures proper package context isolation when an exam exists in multiple packages
func (r *examRepository) GetUserAttemptForExamInPackage(userID uint, examID uint, packageID uint) (*models.UserExamAttempt, error) {
	var attempt models.UserExamAttempt
	err := r.db.Where("user_id = ? AND exam_id = ? AND package_id = ?", userID, examID, packageID).
		Order("created_at DESC").
		First(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No attempt found, not an error
		}
		return nil, fmt.Errorf("failed to fetch user attempt for exam in package: %w", err)
	}
	return &attempt, nil
}

// CreateExamAttemptWithExam creates a new exam attempt using provided exam data
// Optimized to avoid additional DB call to fetch exam details
func (r *examRepository) CreateExamAttemptWithExam(userID uint, exam *models.Exam, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error) {
	// Generate session ID
	sessionID := r.generateSessionID()
	now := time.Now()

	// Create attempt record using provided exam data
	attempt := &models.UserExamAttempt{
		UserID:           userID,
		ExamID:           exam.ID,
		PackageID:        packageID,
		AttemptNumber:    1,
		Status:           models.AttemptStatusStarted,
		StartedAt:        now,
		SessionID:        &sessionID,
		LastActivityAt:   &now,
		TimeLimitSeconds: exam.DurationMinutes * 60,
		TotalQuestions:   exam.TotalQuestions,
		PassingScore:     exam.PassingScore,
		AnswersData:      "{}",
		IsScored:         false,
	}

	if err := r.db.Create(attempt).Error; err != nil {
		return nil, fmt.Errorf("failed to create exam attempt: %w", err)
	}

	return attempt, nil
}

// GetPackageIDForExam retrieves the package ID for an exam via the junction table
func (r *examRepository) GetPackageIDForExam(examID uint) (uint, error) {
	var packageExam models.PackageExam
	err := r.db.Where("exam_id = ? AND is_active = true", examID).First(&packageExam).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrExamNotFound
		}
		return 0, fmt.Errorf("failed to fetch package for exam: %w", err)
	}
	return packageExam.PackageID, nil
}

// MarkExpiredSessionsAsAbandoned marks all expired exam sessions as abandoned
func (r *examRepository) MarkExpiredSessionsAsAbandoned(currentTime time.Time, gracePeriodSeconds int) (int64, error) {
	// Initialize logger with repository context
	log := logger.WithService("ExamRepository").WithFields(logrus.Fields{
		"operation":            "MarkExpiredSessionsAsAbandoned",
		"current_time":         currentTime.Format(time.RFC3339),
		"grace_period_seconds": gracePeriodSeconds,
	})

	log.Debug("Starting bulk update of expired exam sessions")

	// Build the query to find expired sessions
	// A session is considered expired if:
	// 1. Status is STARTED
	// 2. started_at + exam.duration_minutes + grace_period < current_time

	// First, let's get count of sessions that will be updated for logging
	var countQuery = `
		SELECT COUNT(*)
		FROM user_exam_attempts uea
		INNER JOIN exams e ON uea.exam_id = e.id
		WHERE uea.status = ?
		AND uea.started_at + INTERVAL (e.duration_minutes * 60 + ?) SECOND < ?
	`

	var expiredCount int64
	err := r.db.Raw(countQuery, models.AttemptStatusStarted, gracePeriodSeconds, currentTime).Scan(&expiredCount).Error
	if err != nil {
		log.WithError(err).Error("Failed to count expired sessions")
		return 0, fmt.Errorf("failed to count expired sessions: %w", err)
	}

	log.WithField("expired_sessions_found", expiredCount).Debug("Found expired sessions to update")

	if expiredCount == 0 {
		log.Debug("No expired sessions found, skipping update")
		return 0, nil
	}

	// Perform the bulk update
	updateQuery := `
		UPDATE user_exam_attempts uea
		INNER JOIN exams e ON uea.exam_id = e.id
		SET uea.status = ?, uea.updated_at = NOW()
		WHERE uea.status = ?
		AND uea.started_at + INTERVAL (e.duration_minutes * 60 + ?) SECOND < ?
	`

	startTime := time.Now()
	result := r.db.Exec(updateQuery, models.AttemptStatusAbandoned, models.AttemptStatusStarted, gracePeriodSeconds, currentTime)

	if result.Error != nil {
		log.WithError(result.Error).Error("Failed to execute bulk update of expired sessions")
		return 0, fmt.Errorf("failed to mark expired sessions as abandoned: %w", result.Error)
	}

	updatedRows := result.RowsAffected
	duration := time.Since(startTime)

	log.WithFields(logrus.Fields{
		"updated_sessions":     updatedRows,
		"expected_sessions":    expiredCount,
		"update_duration_ms":   duration.Milliseconds(),
		"current_time":         currentTime.Format(time.RFC3339),
		"grace_period_seconds": gracePeriodSeconds,
	}).Info("Successfully marked expired sessions as abandoned")

	// Log warning if actual updated count doesn't match expected count
	if updatedRows != expiredCount {
		log.WithFields(logrus.Fields{
			"expected": expiredCount,
			"actual":   updatedRows,
			"diff":     expiredCount - updatedRows,
		}).Warn("Mismatch between expected and actual updated session count")
	}

	return updatedRows, nil
}

// Custom errors
var (
	ErrExamNotFound         = errors.New("exam not found")
	ErrAttemptNotFound      = errors.New("exam attempt not found")
	ErrExamAlreadySubmitted = errors.New("exam already submitted")
	ErrNotEnrolledInPackage = errors.New("user not enrolled in package")
)
