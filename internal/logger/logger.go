package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Patterns that look like credentials or secrets — redacted in debug logs
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key|api[_-]?secret|auth[_-]?token|access[_-]?token|secret[_-]?key|private[_-]?key|password|passwd|credential)s?\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)--(api-key|password|passwd|token|secret)\s+\S+`),
	regexp.MustCompile(`(?i)(Authorization|X-API-Key)\s*:\s*\S+`),
	regexp.MustCompile(`(?i)(sk-[a-zA-Z0-9]{20,})`),                                    // OpenAI/Anthropic-style keys
	regexp.MustCompile(`(?i)([a-zA-Z0-9+/]{30,}={0,2})`),                                 // base64-like tokens in command context
	regexp.MustCompile(`(?i)curl\s+.*(-u|--user)\s+\S+`),                                 // curl basic auth
	regexp.MustCompile(`(?i)export\s+\S*(API[_-]?KEY|TOKEN|SECRET|PASSWORD)\S*\s*=\s*\S+`),
}

// sanitize scrubs potential secrets from log messages
func sanitize(s string) string {
	for _, pat := range secretPatterns {
		s = pat.ReplaceAllStringFunc(s, func(match string) string {
			// Preserve the key/flag name but redact the value
			if idx := strings.Index(match, "="); idx >= 0 {
				return match[:idx+1] + "[REDACTED]"
			}
			if idx := strings.Index(match, ":"); idx >= 0 {
				return match[:idx+1] + " [REDACTED]"
			}
			// For flag-style: --flag value
			parts := strings.Fields(match)
			if len(parts) >= 2 {
				return parts[0] + " [REDACTED]"
			}
			return "[REDACTED]"
		})
	}
	return s
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	mu         sync.Mutex
	logFile    *os.File
	currentLvl = LevelInfo
)

// Init initializes the logger writing to logFilePath
func Init(logFilePath string, debug bool) {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}

	// get log level from env
	lvlEnv := strings.ToLower(os.Getenv("IRIS_LOG_LEVEL"))
	switch lvlEnv {
	case "debug":
		currentLvl = LevelDebug
	case "info":
		currentLvl = LevelInfo
	case "warn":
		currentLvl = LevelWarn
	case "error":
		currentLvl = LevelError
	default:
		if debug {
			currentLvl = LevelDebug
		} else {
			currentLvl = LevelInfo
		}
	}

	// rotate if log file is larger than 5MB
	if info, err := os.Stat(logFilePath); err == nil && info.Size() > 5*1024*1024 {
		_ = os.Rename(logFilePath, logFilePath+".old")
	}

	_ = os.MkdirAll(filepath.Dir(logFilePath), 0755)
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		logFile = f
	}
}

// Close closes the underlying log file
func Close() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}

func logmsg(lvl Level, lvlStr string, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil || lvl < currentLvl {
		return
	}

	// get caller information to append file:line
	caller := "unknown:0"
	if _, file, line, ok := runtime.Caller(2); ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	tStr := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	msg := fmt.Sprintf(format, a...)
	// Scrub credentials from debug-level messages to prevent secrets leaking into logs
	if lvl == LevelDebug {
		msg = sanitize(msg)
	}
	_, _ = fmt.Fprintf(logFile, "%s [%s] [%s] %s\n", tStr, lvlStr, caller, msg)
}

// Debugf writes a debug log message
func Debugf(format string, a ...any) {
	logmsg(LevelDebug, "DEBUG", format, a...)
}

// Infof writes an info log message
func Infof(format string, a ...any) {
	logmsg(LevelInfo, "INFO", format, a...)
}

// Warnf writes a warning log message
func Warnf(format string, a ...any) {
	logmsg(LevelWarn, "WARN", format, a...)
}

// Errorf writes an error log message
func Errorf(format string, a ...any) {
	logmsg(LevelError, "ERROR", format, a...)
}
