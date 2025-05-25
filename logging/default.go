package logging

import (
	"fmt"
	"strings"
	"time"
)

// defaultLogger is a no-frills adapter that mimics Dyffi’s own color-coded request logs.
type defaultLogger struct{}

// NewDefaultLogger returns a logger that prints colorized output to stdout.
func NewDefaultLogger() *defaultLogger {
	return &defaultLogger{}
}

func (l *defaultLogger) LogInfo(msg string, fields ...Field) {
	l.logWithLevel("INFO", msg, fields...)
}

func (l *defaultLogger) LogWarning(msg string, fields ...Field) {
	l.logWithLevel("WARN", msg, fields...)
}

func (l *defaultLogger) LogError(msg string, fields ...Field) {
	l.logWithLevel("ERROR", msg, fields...)
}

func (l *defaultLogger) LogPanic(msg string, fields ...Field) {
	l.logWithLevel("PANIC", msg, fields...)
	panic(msg)
}

// logWithLevel does the heavy lifting: timestamp, colors, fields.
func (l *defaultLogger) logWithLevel(level, msg string, fields ...Field) {
	// timestamp in YYYY/MM/DD HH:MM:SS
	ts := time.Now().Format("2006/01/02 15:04:05")

	// colors
	reset        := "\033[0m"
	tsColor      := "\033[1;31m" // red for timestamp
	levelColor   := getLevelColor(level)
	msgColor     := "\033[1;36m" // cyan for message
	keyColor     := "\033[1;33m" // yellow for field key
	valueColor   := "\033[1;32m" // green for field value

	// build fields string
	var fieldsStr string
	if len(fields) > 0 {
		parts := make([]string, len(fields))
		for i, f := range fields {
			parts[i] = fmt.Sprintf(
				"%s%s%s: %s%v%s",
				keyColor, f.Key, reset,
				valueColor, f.Value, reset,
			)
		}
		fieldsStr = " | " + strings.Join(parts, ", ")
	}

	// print it out
	fmt.Printf(
		"%s%s%s | Level: %s%s%s | %s%s%s%s\n",
		tsColor, ts, reset,
		levelColor, level, reset,
		msgColor, msg, reset,
		fieldsStr,
	)
}

// getLevelColor returns a bright color code for each log level.
func getLevelColor(level string) string {
	switch level {
	case "INFO":
		return "\033[1;32m" // green
	case "WARN":
		return "\033[1;33m" // yellow
	case "ERROR":
		return "\033[1;31m" // red
	case "PANIC":
		return "\033[1;31m" // red
	default:
		return "\033[0m"
	}
}
