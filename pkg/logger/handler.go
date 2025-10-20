package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TraceHandler wraps an slog.Handler and adds trace context + span events
type TraceHandler struct {
	handler slog.Handler
}

// NewTraceHandler creates a new handler that adds trace_id, span_id and span events
func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{handler: h}
}

func (h *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract trace context from the context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)

		// Add log as span event (visible in Jaeger UI)
		span := trace.SpanFromContext(ctx)
		if span.IsRecording() {
			attrs := make([]attribute.KeyValue, 0)
			r.Attrs(func(a slog.Attr) bool {
				// Convert slog attributes to OTel attributes
				switch a.Value.Kind() {
				case slog.KindString:
					attrs = append(attrs, attribute.String(a.Key, a.Value.String()))
				case slog.KindInt64:
					attrs = append(attrs, attribute.Int64(a.Key, a.Value.Int64()))
				case slog.KindFloat64:
					attrs = append(attrs, attribute.Float64(a.Key, a.Value.Float64()))
				case slog.KindBool:
					attrs = append(attrs, attribute.Bool(a.Key, a.Value.Bool()))
				default:
					attrs = append(attrs, attribute.String(a.Key, a.Value.String()))
				}
				return true
			})

			// Add level as attribute
			attrs = append(attrs, attribute.String("level", r.Level.String()))

			// Add event to span
			span.AddEvent(r.Message, trace.WithAttributes(attrs...))

			// If error level, mark span as error
			if r.Level >= slog.LevelError {
				span.SetStatus(codes.Error, r.Message)
			}
		}
	}
	return h.handler.Handle(ctx, r)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{handler: h.handler.WithGroup(name)}
}
