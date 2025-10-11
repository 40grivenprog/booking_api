package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/vention/booking_api/internal/api"
	"github.com/vention/booking_api/internal/api/middleware"
	"github.com/vention/booking_api/internal/config"
	"github.com/vention/booking_api/internal/database"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/token"
)

func Start(ctx context.Context, cfg *config.Config, logger zerolog.Logger) error {
	// Initialize database
	database, err := database.NewPostgreSQL(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close()

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())

	// Set base logger in context first
	r.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	r.Use(middleware.RequestID()) // Use our custom middleware
	r.Use(middleware.Logger())    // Use our combined logger middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize repository
	queries := db.New(database.DB)

	// Initialize JWT token maker
	tokenMaker, err := token.NewJWTMaker(cfg.JWT.Secret)
	if err != nil {
		return fmt.Errorf("failed to create token maker: %w", err)
	}

	// Apply JWT authentication to all /api routes
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(tokenMaker))

	// Register API routes with JWT protection
	if err := api.Register(ctx, cfg, apiGroup, queries); err != nil {
		return fmt.Errorf("failed to register API routes: %w", err)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := database.Health(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	logger.Info().
		Str("address", srv.Addr).
		Msg("Starting HTTP server")

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	logger.Info().Msg("Shutting down HTTP server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("HTTP server forced to shutdown")
		return err
	}

	logger.Info().Msg("HTTP server exited")
	return nil
}
