package mylog

import (
	"log"
	"os"
)
import "context"

// LogWriteConfig is the local file config.
type LogWriteConfig struct {
	LogLevel Level
	// LogPath is the log path like /usr/local/trpc/log/.
	LogPath string `yaml:"log_path"`
	// Filename is the file name like trpc.log.
	Filename string `yaml:"filename"`
	// MaxAge is the max expire times(day).
	MaxAge int `yaml:"max_age"`
	// MaxSize is the max size of log file(MB).
	MaxSize int `yaml:"max_size"`
	// MaxBackups is the max backup files.
	MaxBackups int `yaml:"max_backups"`
}

type Logger struct {
	logCft   LogWriteConfig
	log      *log.Logger
	fileFd   *os.File
	fullPath string
}

func New(cft *LogWriteConfig) *Logger {
	logger := &Logger{
		logCft: LogWriteConfig{
			LogLevel: DebugLevel,
		},
	}
	if cft != nil {
		logger.logCft = *cft
	}
	// 暂时实现按1天切换
	logger.logCft.MaxAge = 1
	logger.Init()
	return logger
}

type Log interface {
	Init()
	LogFileCheck() error
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

