package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	reset  = "\033[0m"
	blue   = "\033[0;34m"
	yellow = "\033[1;33m"
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	purple = "\033[0;35m"
)

type Logger struct {
	quiet bool
	file  *os.File
}

func New(path string, quiet bool) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{quiet: quiet, file: f}, nil
}

func (l *Logger) Close() { l.file.Close() }

func (l *Logger) write(level, color, symbol, msg string) {
	fmt.Fprintf(l.file, "[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, msg)
	if !l.quiet {
		fmt.Printf("%s%s  %s%s\n", color, symbol, msg, reset)
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.write("INFO", blue, "ℹ", fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(format string, args ...any) {
	l.write("WARN", yellow, "⚠", fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...any) {
	l.write("ERROR", red, "✗", fmt.Sprintf(format, args...))
}

func (l *Logger) Success(format string, args ...any) {
	l.write("SUCCESS", green, "✓", fmt.Sprintf(format, args...))
}

func (l *Logger) Progress(format string, args ...any) {
	l.write("PROGRESS", purple, "→", fmt.Sprintf(format, args...))
}
