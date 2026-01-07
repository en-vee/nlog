package nlog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelTrace slog.Level = slog.LevelDebug - 4
	LevelFatal slog.Level = slog.LevelError + 4
)

var levelMap = map[string]slog.Level{
	"TRACE": LevelTrace,
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
	"FATAL": LevelFatal,
}

type LoggerOpt func(*Logger)

func WithFileLogger(fileLogger *lumberjack.Logger) LoggerOpt {
	return func(logger *Logger) {
		logger.fileLogger = fileLogger
		logger.fileEnabled = true
	}
}

func WithConsoleLogger(enabled bool) LoggerOpt {
	return func(logger *Logger) {
		logger.consoleEnabled = enabled
	}
}

func WithLevel(level string) LoggerOpt {
	return func(logger *Logger) {
		if l, ok := levelMap[strings.ToUpper(level)]; ok {
			logger.level.Set(l)
		}
	}
}

func WithTimestampFormat(timestampFormat string) LoggerOpt {
	return func(logger *Logger) {
		logger.replaceAttrFn = makeReplaceAttrFn(timestampFormat)
	}
}

// Logger provides the ability to log to a file and/or stdout
type Logger struct {
	fileLogger                  *lumberjack.Logger
	fileEnabled, consoleEnabled bool
	level                       *slog.LevelVar
	replaceAttrFn               func([]string, slog.Attr) slog.Attr
	//fileLogger                  *slog.Logger
	//consoleLogger               *slog.Logger
	loggers []*slog.Logger
}

func NewLogger(loggerOpts ...LoggerOpt) *Logger {
	logger := &Logger{}
	logger.loggers = make([]*slog.Logger, 0)
	// Set defaults
	logger.replaceAttrFn = makeReplaceAttrFn(time.RFC3339Nano)
	logger.level = &slog.LevelVar{}

	// Override with LoggerOpt
	for _, opt := range loggerOpts {
		opt(logger)
	}

	slogHandlerOptions := &slog.HandlerOptions{Level: logger.level, ReplaceAttr: logger.replaceAttrFn}

	if logger.consoleEnabled {
		logger.loggers = append(logger.loggers, slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions)))
	}

	if logger.fileEnabled {
		logger.loggers = append(logger.loggers, slog.New(slog.NewTextHandler(logger.fileLogger, slogHandlerOptions)))
	}
	return logger
}

func (logger *Logger) Tracef(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelDebug, makeCustomLogFn(lg, LevelTrace), msg, args...)
	}
}

func (logger *Logger) Debugf(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelDebug, lg.Debug, msg, args...)
	}
}

func (logger *Logger) Infof(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelInfo, lg.Info, msg, args...)
	}
}

func (logger *Logger) Warnf(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelWarn, lg.Warn, msg, args...)
	}
}

func (logger *Logger) Errorf(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelError, lg.Error, msg, args...)
	}
}

func (logger *Logger) Fatalf(msg string, args ...any) {
	for _, lg := range logger.loggers {
		logger.logMsg(slog.LevelDebug, makeCustomLogFn(lg, LevelFatal), msg, args...)
	}
}

func makeCustomLogFn(lg *slog.Logger, customLogLevel slog.Level) func(msg string, args ...any) {
	return func(msg string, args ...any) {
		lg.Log(context.Background(), customLogLevel, msg, args...)
	}
}

func (logger *Logger) logMsg(msgLogLevel slog.Level, f func(msg string, args ...any), msg string, args ...any) {
	if logger.level.Level() > msgLogLevel {
		return
	}

	if len(args) == 0 {
		f(msg)
		return
	}

	f(fmt.Sprintf(msg, args...))

}

func makeReplaceAttrFn(timestampFormat string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			switch {
			case level <= LevelTrace:
				a.Value = slog.StringValue("TRACE")
			case level >= LevelFatal:
				a.Value = slog.StringValue("FATAL")
			}
		}
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(time.Now().Format(timestampFormat))
		}
		return a
	}
}
