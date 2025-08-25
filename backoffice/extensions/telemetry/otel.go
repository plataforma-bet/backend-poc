package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	trace.Span
}

func (span Span) HandlerError(err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	slog.Error("unexpected error",
		slog.Any("error", err),
		slog.Any("span_id", span.SpanContext().SpanID().String()),
		slog.Any("trace_id", span.SpanContext().TraceID().String()),
	)
}

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

var GlobalTracer = GetTracer(scope)

func HandlerUnexpectedError(ctx context.Context, err error) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	slog.Error("unexpected error", slog.Any("error", err))
}

func Shutdown(ctx context.Context) error {
	if err := shutdownTracer(ctx); err != nil {
		return err
	}

	if err := shutdownMeterProvider(ctx); err != nil {
		return err
	}

	return nil
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	ctx, span := GlobalTracer.Start(ctx, name, opts...)

	if requestID, ok := RequestIDFromContext(ctx); ok {
		span.SetAttributes(attribute.String("request_id", requestID))
	}

	span.SetStatus(codes.Ok, "started")

	return ctx, Span{span}
}

func Initialize(ctx context.Context, telemetry OpenTelemetry) error {
	r, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(telemetry.ApplicationName),
			semconv.ServiceNamespace(telemetry.Namespace),
			semconv.ServiceInstanceID(telemetry.InstanceID),
			semconv.ServiceVersion(telemetry.Version),
			semconv.DeploymentEnvironment(telemetry.Environment),
		),
	)
	if err != nil {
		return err
	}

	if err = initializeTracerProvider(ctx, telemetry, r); err != nil {
		return err
	}

	if err = initializeMeterProvider(ctx, telemetry, r); err != nil {
		return err
	}

	initializePropagator()

	return nil

}

func initializePropagator() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		b3.New(b3.WithInjectEncoding(b3.B3SingleHeader|b3.B3MultipleHeader)),
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
