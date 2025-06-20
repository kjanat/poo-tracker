package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// App represents the main application
type App struct {
	container *Container
	server    *http.Server
}

// New creates a new application instance
func New() (*App, error) {
	container, err := NewContainer()
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	app := &App{
		container: container,
	}

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	// Setup Gin router
	router := a.setupRouter()

	// Create HTTP server with security and performance timeouts
	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.container.Config.Host, a.container.Config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second, // Maximum duration for reading the entire request
		WriteTimeout: 15 * time.Second, // Maximum duration for writing the response
		IdleTimeout:  60 * time.Second, // Maximum amount of time to wait for the next request when keep-alives are enabled
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s:%s", a.container.Config.Host, a.container.Config.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	// Cleanup resources
	if err := a.container.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup resources: %w", err)
	}

	log.Println("Server exited")
	return nil
}

// setupRouter configures the Gin router
func (a *App) setupRouter() *gin.Engine {
	// Set Gin mode based on environment
	if a.container.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API routes will be added in later phases
	// For now, just add a placeholder
	api := router.Group("/api/v1")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API is running",
				"version": "v1",
			})
		})
	}

	return router
}
