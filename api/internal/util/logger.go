package util

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime"
	"time"

	"github.com/fatih/color"
)

func init() {
	// Force colorization even when output is not a terminal
	color.NoColor = false
}

type Logger struct {
	Log *slog.Logger
}

func colorizeStackTrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	stackTrace := ""
	for {
		frame, more := frames.Next()

		stackTrace += color.BlueString(frame.Function) + "\n"
		stackTrace += "\t" + color.GreenString(frame.File)
		stackTrace += ":" + color.YellowString(
			fmt.Sprintf("%d", frame.Line),
		) + "\n"

		if !more {
			break
		}
	}
	return stackTrace + "\n"
}

func (l *Logger) Err(msg string, args ...any) {
	l.Log.Error(msg, args...)
	log.Println(colorizeStackTrace())
}

func (l *Logger) Dbg(msg string, args ...any) {
	// l.Log.Debug(msg, args...)
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(
		time.Now(),
		slog.LevelDebug,
		fmt.Sprintf(msg, args...),
		pcs[0],
	)
	_ = l.Log.Handler().Handle(context.Background(), r)
}

func (l *Logger) Inf(msg string, args ...any) {
	l.Log.Info(msg, args...)
}
