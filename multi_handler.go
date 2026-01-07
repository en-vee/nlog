package nlog

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type HandlerOption func(*MultiHandler)

func WithFile(fileLogger *lumberjack.Logger) HandlerOption {
	return func(h *MultiHandler) {
		h.handlers = append(h.handlers, slog.NewTextHandler(fileLogger, &slog.HandlerOptions{Level: h.level.Level(), ReplaceAttr: h.replaceAttrFn}))
	}
}

func WithConsole() HandlerOption {
	return func(h *MultiHandler) {
		h.handlers = append(h.handlers, slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: h.level.Level(), ReplaceAttr: h.replaceAttrFn}))
	}
}

func WithLogLevel(level string) HandlerOption {
	return func(h *MultiHandler) {
		if l, ok := levelMap[strings.ToUpper(level)]; ok {
			h.level.Set(l)
		}
	}
}

func WithLogTimestampFormat(timestampFormat string) HandlerOption {
	return func(h *MultiHandler) {
		h.replaceAttrFn = makeReplaceAttrFn(timestampFormat)
	}
}

type MultiHandler struct {
	handlers      []slog.Handler
	level         slog.LevelVar
	replaceAttrFn func([]string, slog.Attr) slog.Attr
}

func NewMultiHandler(opts ...HandlerOption) *MultiHandler {
	h := &MultiHandler{handlers: make([]slog.Handler, 0), replaceAttrFn: makeReplaceAttrFn(time.RFC3339Nano)}
	// Set defaults

	// Override
	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {

	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	mh := &MultiHandler{}
	mh.handlers = make([]slog.Handler, 0)
	for _, handler := range h.handlers {
		mh.handlers = append(mh.handlers, handler.WithAttrs(attrs))
	}
	return mh
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	mh := &MultiHandler{}
	mh.handlers = make([]slog.Handler, 0)
	for _, handler := range h.handlers {
		mh.handlers = append(mh.handlers, handler.WithGroup(name))
	}
	return mh
}
