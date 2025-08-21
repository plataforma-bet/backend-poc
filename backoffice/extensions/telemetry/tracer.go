package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initializeTracerProvider(ctx context.Context, telemetry OpenTelemetry, resource *resource.Resource) error {
	tracerExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(telemetry.Endpoint),
		otlptracegrpc.WithInsecure(),
	))
	if err != nil {
		return err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(telemetry.SamplingRate)),
		sdktrace.WithBatcher(tracerExporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)

	return nil
}

func shutdownTracer(ctx context.Context) error {
	provider, ok := any(otel.GetTracerProvider()).(*sdktrace.TracerProvider)
	if !ok {
		return fmt.Errorf("otel.GetTracerProvider() isn't set or it's invalid")
	}

	if err := provider.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
