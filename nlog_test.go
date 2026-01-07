package nlog

import (
	"testing"

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
