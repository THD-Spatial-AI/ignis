package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var (
	Info  *log.Logger
	Debug *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	// Current log level - can be configured
	currentLogLevel LogLevel = LogLevelInfo // Default to INFO level
)

// InitLogger initializes loggers with configurable debug level
func InitLogger() {
	// Set log level from environment variable (with fallback)
	setLogLevelFromEnv()

	// Try to create log directory, but don't fail if it can't be created
	logDir := "logs"
	var logFile *os.File

	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If we can't create log directory, just use console logging
		log.Printf("Warning: Failed to create log directory, using console only: %v", err)
	} else {
		// Try to open log file
		logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
		if file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log.Printf("Warning: Failed to open log file, using console only: %v", err)
		} else {
			logFile = file
		}
	}

	// Create multi-writers for console and file output
	var multiInfo, multiWarn, multiError, multiDebug io.Writer

	if logFile != nil {
		multiInfo = io.MultiWriter(os.Stdout, logFile)
		multiWarn = io.MultiWriter(os.Stdout, logFile)
		multiError = io.MultiWriter(os.Stderr, logFile) // Errors to stderr

		// Debug output depends on log level
		if currentLogLevel <= LogLevelDebug {
			multiDebug = io.MultiWriter(os.Stdout, logFile)
		} else {
			multiDebug = logFile // Debug only to file when not in debug mode
		}
	} else {
		// Console-only logging if file creation failed
		multiInfo = os.Stdout
		multiWarn = os.Stdout
		multiError = os.Stderr
		multiDebug = os.Stdout
	}

	// Initialize loggers
	Info = log.New(multiInfo, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(multiDebug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(multiWarn, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Log the current log level
	Info.Printf("Logger initialized with level: %s", getLogLevelName(currentLogLevel))
}

// setLogLevelFromEnv sets the log level from LOG_LEVEL environment variable
func setLogLevelFromEnv() {
	// Use os.Getenv directly to avoid circular dependency issues
	logLevelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevelStr == "" {
		logLevelStr = "INFO" // Default value
	}

	switch logLevelStr {
	case "DEBUG":
		currentLogLevel = LogLevelDebug
	case "INFO":
		currentLogLevel = LogLevelInfo
	case "WARN", "WARNING":
		currentLogLevel = LogLevelWarn
	case "ERROR":
		currentLogLevel = LogLevelError
	default:
		currentLogLevel = LogLevelInfo
		log.Printf("Unknown log level '%s', defaulting to INFO", logLevelStr)
	}
}

// getLogLevelName returns the name of the log level
func getLogLevelName(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// SetLogLevel allows programmatic control of log level
func SetLogLevel(level LogLevel) {
	currentLogLevel = level
	Info.Printf("Log level changed to: %s", getLogLevelName(level))
}

// IsDebugEnabled returns true if debug logging is enabled
func IsDebugEnabled() bool {
	return currentLogLevel <= LogLevelDebug
}

// GetLogLevel returns the current log level
func GetLogLevel() LogLevel {
	return currentLogLevel
}
