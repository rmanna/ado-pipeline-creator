// Package logger reference: https://www.datadoghq.com/blog/go-logging/
package logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// Log exported
var Log *CustomLogger

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// CustomLogger enforces specific log message formats
type CustomLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() {

	var baseLogger = logrus.New()
	Log = &CustomLogger{baseLogger}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/app.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	Log.SetLevel(logrus.InfoLevel)
	Log.SetOutput(colorable.NewColorableStdout())
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	Log.AddHook(rotateFileHook)
}

// Declare variables to store log messages as new Events
var (
	infoArg                = Event{1, "%s"}
	invalidArgMessage      = Event{2, "Invalid arg: %s"}
	invalidArgValueMessage = Event{3, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{4, "Missing arg: %s"}
)

// RequestFields is a standard HTTP request message
func (l *CustomLogger) RequestFields(methodValue string, pathValue string) {
	l.WithFields(logrus.Fields{
		"method":     methodValue,
		"path":       pathValue,
		"request_id": uuid.New(),
	}).Infof(infoArg.message, "Request")
}

// InfoArg is a standard info message
func (l *CustomLogger) InfoArg(argumentName string) {
	l.Infof(infoArg.message, argumentName)
}

// InvalidArg is a standard error message
func (l *CustomLogger) InvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func (l *CustomLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (l *CustomLogger) MissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}
