package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Log struct {
	logger      *slog.Logger
	packageName string
}

func NewLogger(packageName string) *Log {
	logLevel := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))

	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	lvl := new(slog.LevelVar)
	lvl.Set(level)
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: lvl,
		}))

	return &Log{logger: logger, packageName: packageName}
}

func (l Log) Debug(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	l.logger.Debug(log, "package", l.packageName)
}

func (l Log) Info(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	l.logger.Info(log, "package", l.packageName)
}

func (l Log) Warn(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	l.logger.Warn(log, "package", l.packageName)
}

func (l Log) Error(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	l.logger.Error(log, "package", l.packageName)
}
