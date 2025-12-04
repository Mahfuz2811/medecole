package service

import (
	"context"
	"github.com/Mahfuz2811/medecole/backend/internal/logger"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"time"

	"github.com/sirupsen/logrus"
)

// ExamCleanupService handles background cleanup of expired exam sessions
type ExamCleanupService interface {
	Start(ctx context.Context) error
	Stop() error
}

// examCleanupService implements ExamCleanupService
type examCleanupService struct {
	examRepo        repository.ExamRepository
	cleanupInterval time.Duration
	gracePeriod     time.Duration
	stopChan        chan struct{}
	stopped         bool
}

// CleanupConfig holds configuration for the cleanup service
type CleanupConfig struct {
	CleanupInterval time.Duration // How often to run cleanup (default: 1 minute)
	GracePeriod     time.Duration // Grace period to avoid race conditions (default: 2 minutes)
}

// NewExamCleanupService creates a new exam cleanup service
func NewExamCleanupService(examRepo repository.ExamRepository, config CleanupConfig) ExamCleanupService {
	// Set default values if not provided
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 1 * time.Minute
	}
	if config.GracePeriod == 0 {
		config.GracePeriod = 2 * time.Minute
	}

	return &examCleanupService{
		examRepo:        examRepo,
		cleanupInterval: config.CleanupInterval,
		gracePeriod:     config.GracePeriod,
		stopChan:        make(chan struct{}),
		stopped:         false,
	}
}

// Start begins the background cleanup process
func (s *examCleanupService) Start(ctx context.Context) error {
	log := logger.WithService("ExamCleanupService").WithFields(logrus.Fields{
		"operation":        "Start",
		"cleanup_interval": s.cleanupInterval.String(),
		"grace_period":     s.gracePeriod.String(),
	})

	log.Info("Starting exam cleanup background service")

	// Create a ticker for periodic cleanup
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	// Run initial cleanup
	s.performCleanup()

	for {
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, stopping exam cleanup service")
			return ctx.Err()
		case <-s.stopChan:
			log.Info("Stop signal received, stopping exam cleanup service")
			return nil
		case <-ticker.C:
			s.performCleanup()
		}
	}
}

// Stop stops the cleanup service
func (s *examCleanupService) Stop() error {
	if s.stopped {
		return nil
	}

	log := logger.WithService("ExamCleanupService").WithField("operation", "Stop")
	log.Info("Stopping exam cleanup service")

	s.stopped = true
	close(s.stopChan)
	return nil
}

// performCleanup performs the actual cleanup of expired sessions
func (s *examCleanupService) performCleanup() {
	startTime := time.Now()
	log := logger.WithService("ExamCleanupService").WithFields(logrus.Fields{
		"operation":    "PerformCleanup",
		"started_at":   startTime.Format(time.RFC3339),
		"grace_period": s.gracePeriod.String(),
	})

	log.Debug("Starting cleanup of expired exam sessions")

	// Use current time for accurate expiration calculation
	currentTime := time.Now()
	gracePeriodSeconds := int(s.gracePeriod.Seconds())

	log.WithFields(logrus.Fields{
		"current_time":         currentTime.Format(time.RFC3339),
		"grace_period_seconds": gracePeriodSeconds,
	}).Debug("Calculated parameters for expired session cleanup")

	// Find and update expired sessions
	updatedCount, err := s.examRepo.MarkExpiredSessionsAsAbandoned(currentTime, gracePeriodSeconds)
	if err != nil {
		log.WithError(err).Error("Failed to mark expired sessions as abandoned")
		return
	}

	// Calculate performance metrics
	duration := time.Since(startTime)
	log.WithFields(logrus.Fields{
		"updated_sessions":     updatedCount,
		"cleanup_duration":     duration.String(),
		"cleanup_duration_ms":  duration.Milliseconds(),
		"current_time":         currentTime.Format(time.RFC3339),
		"grace_period_seconds": gracePeriodSeconds,
	}).Info("Completed cleanup of expired exam sessions")

	// Log warning if cleanup took too long
	if duration > 30*time.Second {
		log.WithFields(logrus.Fields{
			"duration_seconds": duration.Seconds(),
			"threshold":        30,
		}).Warn("Cleanup operation took longer than expected")
	}
}
