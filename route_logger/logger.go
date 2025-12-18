package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Global logger instance
	globalLogger *zap.Logger
)

// InitLogger initializes the global logger based on environment and log level
// env: "dev" for development mode, "prod" for production mode
// level: "debug", "info", "warn", "error", "fatal", "panic"
func InitLogger(env string, level string) (*zap.Logger, error) {
	var config zap.Config
	var err error
	// Parse log level
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	if env == "dev" || env == "development" {
		// Development configuration - human readable console output
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		globalLogger, err = config.Build()
		if err != nil {
			return nil, err
		}
	} else {
		// Production configuration - JSON structured logging with file rotation
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		// its use zapcore for Encoder
		// config := zap.NewProductionEncoderConfig()
		// config.TimeKey = "timestamp"
		// config.EncodeTime = zapcore.ISO8601TimeEncoder
		// config.EncodeLevel = zapcore.LowercaseLevelEncoder
		// config.EncodeCaller = zapcore.ShortCallerEncoder

		// Create logs directory if it doesn't exist
		logsDir := "logs"
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return nil, err
		}

		// Configure lumberjack for log rotation
		lumberjackLogger := &lumberjack.Logger{
			Filename:   filepath.Join(logsDir, "app.log"),
			MaxSize:    50,   // 50 MB
			MaxBackups: 5,    // Keep 5 backup files
			MaxAge:     30,   // Keep logs for 30 days
			Compress:   true, // Compress rotated files
		}

		// Create file writer with rotation
		fileWriter := zapcore.AddSync(lumberjackLogger)

		// Create console writer for production (optional - can be removed if not needed)
		consoleWriter := zapcore.AddSync(os.Stdout)

		// Create encoder
		encoder := zapcore.NewJSONEncoder(config.EncoderConfig)

		// Create core with both file and console output
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, fileWriter, logLevel),
			zapcore.NewCore(encoder, consoleWriter, logLevel),
		)

		// Create logger with the core
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		// wraperLogger = zap.New(core, zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return globalLogger, nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Fallback to development logger if not initialized
		globalLogger, _ = zap.NewDevelopment()
	}
	return globalLogger
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// SetGlobalLogger sets the global logger instance (useful for testing)
func SetGlobalLogger(logger *zap.Logger) {
	globalLogger = logger
}

/*
Example usage:

// Initialize logger
logger, err := InitLogger("dev", "info")
if err != nil {
    log.Fatal("Failed to initialize logger:", err)
}
defer logger.Sync()

// Use logger
logger.Info("Application started", zap.String("version", "1.0.0"))
logger.Error("Failed to connect to database", zap.Error(err))
logger.Debug("Processing request", zap.String("request_id", "123"))

// Or use global logger
globalLogger := GetLogger()
globalLogger.Info("Something happened", zap.String("request_id", "123"))
globalLogger.Error("Failed to connect", zap.Error(err))

// Common log patterns:
// Info logging
logger.Info("User logged in",
    zap.String("user_id", "123"),
    zap.String("email", "user@example.com"),
    zap.Duration("duration", time.Since(start)))

// Error logging
logger.Error("Database connection failed",
    zap.Error(err),
    zap.String("host", "localhost"),
    zap.Int("port", 5432))

// Debug logging
logger.Debug("Processing request",
    zap.String("method", "POST"),
    zap.String("path", "/api/users"),
    zap.String("request_id", "abc123"))

// Warning logging
logger.Warn("Rate limit exceeded",
    zap.String("ip", "192.168.1.1"),
    zap.Int("requests", 100))

// Fatal logging (will exit the program)
logger.Fatal("Critical system failure", zap.Error(err))

// Panic logging (will panic)
logger.Panic("Unrecoverable error", zap.Error(err))
*/
