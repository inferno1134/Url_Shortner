package logger

import (
	"log"
	"os"
	"strings"
)

// Log levels for the application.
// LevelDBOnly prints only database timing logs.
// LevelInfo prints DB timing logs plus regular informational logs.
// LevelDebug prints all logs, including detailed debug output.
type Level int

const (
	LevelDBOnly Level = iota + 1
	LevelInfo
	LevelDebug
)

var currentLevel = LevelDBOnly

func init() {
	SetLevel(ParseLevel(os.Getenv("LOG_LEVEL")))
}

// SetLevel configures the global logger level.
func SetLevel(level Level) {
	currentLevel = level
}

// ParseLevel converts a string value into a log level.
// Supported values: DB_ONLY, DB, DEBUG, INFO.
func ParseLevel(value string) Level {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "DB_ONLY", "DBONLY", "DB":
		return LevelDBOnly
	case "DEBUG":
		return LevelDebug
	default:
		return LevelDBOnly
	}
}

// DB logs are for database timing and lookup measurements only.
func DB(format string, v ...any) {
	if currentLevel >= LevelDBOnly {
		log.Printf("[DB] "+format, v...)
	}
}

// Info logs are for normal operational information.
func Info(format string, v ...any) {
	if currentLevel >= LevelInfo {
		log.Printf(format, v...)
	}
}

// Debug logs are for verbose debug output.
func Debug(format string, v ...any) {
	if currentLevel >= LevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Error logs should always be printed.
func Error(format string, v ...any) {
	log.Printf("[ERROR] "+format, v...)
}
