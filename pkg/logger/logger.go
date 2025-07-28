package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	file        *os.File
)

// Initialize creates the log file and sets up the loggers
func Initialize() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Open log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("logs/app-%s.log", timestamp)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	file = f

	// Create multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Initialize loggers with prefix and flags
	InfoLogger = log.New(multiWriter,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLogger = log.New(multiWriter,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Close closes the log file
func Close() error {
	if file != nil {
		return file.Close()
	}
	return nil
}

// Info logs an info message with optional fields
func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error logs an error message with optional fields
func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// RequestLog represents a structured log for HTTP requests
type RequestLog struct {
	Method     string
	Path       string
	Duration   time.Duration
	StatusCode int
	UserID     string
	Error      error
}

// LogRequest logs HTTP request information
func LogRequest(rl RequestLog) {
	if rl.Error != nil {
		Error("Request: %s %s | Duration: %v | Status: %d | UserID: %s | Error: %v",
			rl.Method, rl.Path, rl.Duration, rl.StatusCode, rl.UserID, rl.Error)
		return
	}

	Info("Request: %s %s | Duration: %v | Status: %d | UserID: %s",
		rl.Method, rl.Path, rl.Duration, rl.StatusCode, rl.UserID)
}
