package util

import (
	"log/slog"
	"runtime/debug"
)

type Logger struct {
	Log *slog.Logger
}

func (l *Logger) Err(msg string, args ...any) {
	l.Log.Error(msg, append(args, "stacktrace", string(debug.Stack()))...)
}

func (l *Logger) Dbg(msg string, args ...any) {
	l.Log.Debug(msg, args...)
}

func (l *Logger) Inf(msg string, args ...any) {
	l.Log.Info(msg, args...)
}
