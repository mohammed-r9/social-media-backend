package logging

import (
	"io"
	"log/slog"
	"os"
)

type Format int

const (
	JSON Format = iota
	Text
)

type LoggerConfig struct {
	Level  slog.Level
	Format Format
	Output io.Writer
	Attrs  []slog.Attr
}

func NewLogger(cfg LoggerConfig) *slog.Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	var handler slog.Handler

	switch cfg.Format {
	case Text:
		handler = slog.NewTextHandler(cfg.Output, opts)
	default:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	logger := slog.New(handler)

	if len(cfg.Attrs) > 0 {
		args := make([]any, 0, len(cfg.Attrs)*2)
		for _, a := range cfg.Attrs {
			args = append(args, a.Key, a.Value.Any())
		}
		logger = logger.With(args...)
	}

	return logger
}
