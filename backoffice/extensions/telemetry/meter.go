package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func initializeMeterProvider(ctx context.Context, telemetry OpenTelemetry, resource *resource.Resource) error {
	metricExplorer, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(telemetry.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExplorer,
				sdkmetric.WithInterval(3*time.Second),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return nil
}

func shutdownMeterProvider(ctx context.Context) error {
	provider, ok := any(otel.GetMeterProvider()).(*sdkmetric.MeterProvider)
	if !ok {
		return fmt.Errorf("otel.GetMeterProvider() is not set or *sdkmetric.MeterProvider")
	}

	if err := provider.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
