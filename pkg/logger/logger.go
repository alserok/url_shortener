package logger

import "context"

type Logger interface {
	Info(msg string, args ...Arg)
	Debug(msg string, args ...Arg)
	Error(msg string, args ...Arg)
	Warn(msg string, args ...Arg)

	Close() error
}

const (
	Slog = iota
)

func New(t uint, env string) Logger {
	switch t {
	case Slog:
		return newSlog(env)
	default:
		panic("invalid logger type")
	}
}

type Arg struct {
	Key string
	Val any
}

func WithArg(key string, val any) Arg {
	return Arg{
		Key: key,
		Val: val,
	}
}

type ContextKey string

const (
	ctxKey ContextKey = "logger_ctx_key"
)

func WrapLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, ctxKey, log)
}

func ExtractLogger(ctx context.Context) Logger {
	return ctx.Value(ctxKey).(Logger)
}
