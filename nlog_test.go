package nlog

import (
	"log/slog"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestXxx(t *testing.T) {
	fileLogger := &lumberjack.Logger{
		Filename: "/tmp/n.log",
	}
	l := NewLogger(WithConsoleLogger(true), WithFileLogger(fileLogger), WithLevel("TRACE"))
	l.Tracef("hello, nlog!")
	l.Fatalf("hello, nlog!")
}

func TestHandler(t *testing.T) {

	fileLogger := &lumberjack.Logger{
		Filename: "/tmp/nmh.log",
	}

	//mh := NewMultiHandler(WithConsole(), WithLogLevel("DEBUG"))
	mh := NewMultiHandler(WithLogTimestampFormat(time.RFC3339Nano), WithFile(fileLogger), WithConsole(), WithLogLevel("INFO"))
	slog.SetDefault(slog.New(mh))
	l := slog.Default().With("interfaceName", "data")

	l.Debug("hello, nlog multi-handler", "sessionId", 1234)
	l.Info("hello, nlog multi-handler", "sessionId", 1234)
}
