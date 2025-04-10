package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger initializes a new zap logger based on LOG_FORMAT environment variable
func NewLogger() *zap.Logger {
	logFormat := os.Getenv("LOG_FORMAT") // Retrieve log format from environment variable

	var config zap.Config
	if logFormat == "json" {
		config = zap.NewProductionConfig() // JSON logging for production
	} else if logFormat == "color" {
		config = zap.NewDevelopmentConfig()                                                    // Colorful console logs
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0") // Custom time format
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                    // Enable ANSI colors
	} else {
		config = zap.NewDevelopmentConfig() // Default to development mode (non-colored logs)
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0")
	}

	logger, _ := config.Build()
	return logger
}
