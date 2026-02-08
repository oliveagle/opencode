package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Config struct {
	Print bool
	Dev   bool
	Level Level
}

type Logger struct {
	service string
	print   bool
	dev     bool
	level   Level
	mu      sync.Mutex
	file    *os.File
}

var Default *Logger
var logFile *os.File
var logMu sync.Mutex

func Init(cfg Config) {
	logMu.Lock()
	defer logMu.Unlock()

	level := cfg.Level
	if level == "" {
		level = LevelInfo
	}

	// Create log directory if needed
	logDir := filepath.Join(os.TempDir(), "opencode", "logs")
	os.MkdirAll(logDir, 0755)

	// Create log file
	logPath := filepath.Join(logDir, fmt.Sprintf("opencode-%s.log", time.Now().Format("2006-01-02")))
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
	}

	Default = &Logger{
		service: "main",
		print:   cfg.Print,
		dev:     cfg.Dev,
		level:   level,
		file:    logFile,
	}
}

func (l *Logger) create(service string) *Logger {
	return &Logger{
		service: service,
		print:   l.print,
		dev:     l.dev,
		level:   l.level,
		file:    l.file,
	}
}

func Create(cfg map[string]interface{}) *Logger {
	service, _ := cfg["service"].(string)
	if service == "" {
		service = "unknown"
	}
	if Default == nil {
		return &Logger{service: service}
	}
	return Default.create(service)
}

func (l *Logger) log(level Level, msg string, data map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339Nano),
		"level":     level,
		"service":   l.service,
		"message":   msg,
	}
	for k, v := range data {
		entry[k] = v
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Write to file
	if l.file != nil {
		jsonData, _ := json.Marshal(entry)
		l.file.WriteString(string(jsonData) + "\n")
	}

	// Print to stderr if enabled
	if l.print {
		if l.dev {
			fmt.Fprintf(os.Stderr, "[%s] %s: %s %+v\n", level, l.service, msg, data)
		} else {
			jsonData, _ := json.Marshal(entry)
			fmt.Fprintln(os.Stderr, string(jsonData))
		}
	}
}

func (l *Logger) shouldLog(level Level) bool {
	levels := map[Level]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) Debug(msg string, data map[string]interface{}) {
	l.log(LevelDebug, msg, data)
}

func (l *Logger) Info(msg string, data map[string]interface{}) {
	l.log(LevelInfo, msg, data)
}

func (l *Logger) Warn(msg string, data map[string]interface{}) {
	l.log(LevelWarn, msg, data)
}

func (l *Logger) Error(msg string, data map[string]interface{}) {
	l.log(LevelError, msg, data)
}

func File() string {
	if logFile != nil {
		return logFile.Name()
	}
	return ""
}
