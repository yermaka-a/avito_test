package logger

import (
	"io"
	"log/slog"
)

func Setup(out io.Writer, level slog.Leveler) *slog.Logger {
	logger := slog.New(
		slog.NewTextHandler(out,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			}),
	)
	return logger
}
