package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/db"
	"github.com/ilbagatto/tarot-api/internal/logging"
	"github.com/ilbagatto/tarot-api/internal/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	logger := logging.NewLogger()
	defer logger.Sync()

	logger.Info("Starting server...")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
	}
	defer database.Close()

	// Initialize application
	application := app.NewApp(database)
	// Add global middleware for charset=utf-8 in JSON responses
	application.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Установим заголовок ЗАРАНЕЕ
			c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
			return next(c)
		}
	})

	routes.InitRoutes(application)

	// Start the server in a goroutine
	go func() {
		if err := application.Echo.Start(":" + port); err != nil {
			logger.Fatal("Shutting down the server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	logger.Info("Shutting down server...", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Echo.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
