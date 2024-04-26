package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	Log zerolog.Logger
}

func NewLogger() *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf(" %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	return &Logger{
		Log: log,
	}
}

func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	logger.Log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	logger.Log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
