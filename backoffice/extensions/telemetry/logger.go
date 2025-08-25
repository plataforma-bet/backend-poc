package telemetry

import (
	cfg "backend-poc/backoffice/config"
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

func SetLogger(cfg cfg.Global) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	h := &Logger{}

	if cfg.Environment == "dev" {
		h.Handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		h.Handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(h).With(slog.Any("dd.env", cfg.Environment))

	slog.SetDefault(logger)
}

type Logger struct {
	Handler slog.Handler
}

func (h *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Logger) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0)

	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		traceID := sc.TraceID().String()
		spanID := sc.SpanID().String()

		attrs = append(attrs,
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)
	}

	if requestID, ok := RequestIDFromContext(ctx); ok {
		attrs = append(attrs, slog.String("request_id", requestID))
	}

	err := h.Handler.WithAttrs(attrs).Handle(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

func (h *Logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.Handler.WithAttrs(attrs)
}

func (h *Logger) WithGroup(name string) slog.Handler {
	return h.Handler.WithGroup(name)
}
