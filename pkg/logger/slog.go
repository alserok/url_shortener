package logger

import (
	"io"
	"log/slog"
	"os"
)

func newSlog(env string) *slogLogger {
	var l *slog.Logger

	switch env {
	case "DEV":
		l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "PROD":
		f, err := os.OpenFile("logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic("failed to open file: " + err.Error())
		}

		l = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		panic("invalid logger type")
	}

	return &slogLogger{
		l: l,
	}
}

type slogLogger struct {
	l *slog.Logger

	f io.WriteCloser
}

func (l *slogLogger) Close() error {
	if l.f != nil {
		return l.f.Close()
	}

	return nil
}

func (l *slogLogger) Info(msg string, args ...Arg) {
	if len(args) == 0 {
		l.l.Info(msg)
	} else {
		vals := make([]any, 0, len(args))
		for _, arg := range args {
			vals = append(vals, slog.Attr{
				Key:   arg.Key,
				Value: slog.AnyValue(arg.Val),
			})
		}

		l.l.Info(msg, vals...)
	}
}

func (l *slogLogger) Debug(msg string, args ...Arg) {
	if len(args) == 0 {
		l.l.Debug(msg)
	} else {
		vals := make([]any, 0, len(args))
		for _, arg := range args {
			vals = append(vals, slog.Attr{
				Key:   arg.Key,
				Value: slog.AnyValue(arg.Val),
			})
		}

		l.l.Debug(msg, vals...)
	}
}

func (l *slogLogger) Error(msg string, args ...Arg) {
	if len(args) == 0 {
		l.l.Error(msg)
	} else {
		vals := make([]any, 0, len(args))
		for _, arg := range args {
			vals = append(vals, slog.Attr{
				Key:   arg.Key,
				Value: slog.AnyValue(arg.Val),
			})
		}

		l.l.Error(msg, vals...)
	}
}

func (l *slogLogger) Warn(msg string, args ...Arg) {
	if len(args) == 0 {
		l.l.Warn(msg)
	} else {
		vals := make([]any, 0, len(args))
		for _, arg := range args {
			vals = append(vals, slog.Attr{
				Key:   arg.Key,
				Value: slog.AnyValue(arg.Val),
			})
		}

		l.l.Warn(msg, vals...)
	}
}
