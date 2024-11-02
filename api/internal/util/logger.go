package util

type Logger interface {
	Err(msg string, args ...any)
	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
}
