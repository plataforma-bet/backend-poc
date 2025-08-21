package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
)

func Initialize(ctx context.Context, telemetry OpenTelemetry) error {
	r, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(),
	)
	if err != nil {
		return err
	}

	if err := initializeTracerProvider(ctx, telemetry, r); err != nil {
		return err
	}

	return nil

}

func Shutdown(ctx context.Context) error {
	if err := shutdownTracer(ctx); err != nil {
		return err
	}

	return nil
}
