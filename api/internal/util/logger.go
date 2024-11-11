package util

import (
	"log"
	"log/slog"
	"runtime/debug"
)

type Logger struct {
	Log *slog.Logger
}

func (l *Logger) Err(msg string, args ...any) {
	l.Log.Error(msg, args...)
	debug.PrintStack()
	log.Println("========")
}

func (l *Logger) Dbg(msg string, args ...any) {
	l.Log.Debug(msg, args...)
}

func (l *Logger) Inf(msg string, args ...any) {
	l.Log.Info(msg, args...)
}
