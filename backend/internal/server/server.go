package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"github.com/Mahfuz2811/medecole/backend/internal/config"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/routes"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// BackgroundServices holds all background services
type BackgroundServices struct {
	cleanupService service.ExamCleanupService
	ctx            context.Context
	cancel         context.CancelFunc
}

// Server represents the HTTP server with background services
type Server struct {
	httpServer         *http.Server
	backgroundServices *BackgroundServices
	db                 *database.Database
}

// NewServer creates a new server instance with all dependencies
func NewServer(cfg *config.Config, db *database.Database, router *gin.Engine) *Server {
	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Initialize background services
	backgroundServices := initializeBackgroundServices(cfg, db)

	return &Server{
		httpServer:         httpServer,
		backgroundServices: backgroundServices,
		db:                 db,
	}
}

// Start starts the server and all background services
func (s *Server) Start() error {
	// Start background services
	if s.backgroundServices != nil {
		s.startBackgroundServices()
	}

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	return s.Shutdown()
}

// Shutdown gracefully shuts down the server and all background services
func (s *Server) Shutdown() error {
	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop background services first
	if s.backgroundServices != nil {
		s.stopBackgroundServices()
	}

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	// Close database connection
	if s.db != nil {
		s.db.Close()
	}

	log.Println("Server exited gracefully")
	return nil
}

// initializeBackgroundServices creates and configures all background services
func initializeBackgroundServices(cfg *config.Config, db *database.Database) *BackgroundServices {
	if !cfg.Cleanup.Enabled {
		log.Println("Background services disabled by configuration")
		return nil
	}

	// Create context for background services
	ctx, cancel := context.WithCancel(context.Background())

	// Create exam repository for cleanup service
	examRepo := routes.CreateExamRepository(db, cfg)

	// Create cleanup service with configuration
	cleanupConfig := service.CleanupConfig{
		CleanupInterval: cfg.Cleanup.CleanupInterval,
		GracePeriod:     cfg.Cleanup.GracePeriod,
	}
	cleanupService := service.NewExamCleanupService(examRepo, cleanupConfig)

	return &BackgroundServices{
		cleanupService: cleanupService,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// startBackgroundServices starts all configured background services
func (s *Server) startBackgroundServices() {
	if s.backgroundServices == nil {
		return
	}

	// Start cleanup service
	if s.backgroundServices.cleanupService != nil {
		go func() {
			log.Printf("Starting exam cleanup service")
			if err := s.backgroundServices.cleanupService.Start(s.backgroundServices.ctx); err != nil && err != context.Canceled {
				log.Printf("Cleanup service error: %v", err)
			}
		}()
	}

	log.Println("All background services started successfully")
}

// stopBackgroundServices gracefully stops all background services
func (s *Server) stopBackgroundServices() {
	if s.backgroundServices == nil {
		return
	}

	log.Println("Stopping background services...")

	// Cancel context to stop all services
	s.backgroundServices.cancel()

	// Stop cleanup service explicitly
	if s.backgroundServices.cleanupService != nil {
		if err := s.backgroundServices.cleanupService.Stop(); err != nil {
			log.Printf("Error stopping cleanup service: %v", err)
		}
	}

	log.Println("All background services stopped successfully")
}
