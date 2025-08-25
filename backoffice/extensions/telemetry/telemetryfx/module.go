package telemetryfx

import (
	"backend-poc/backoffice/extensions/telemetry"
	"context"

	"go.uber.org/fx"
)

func Module(applicationName string) fx.Option {
	return fx.Module("telemetryfx:oTel",
		fx.Provide(func() telemetry.OpenTelemetry {
			return telemetry.NewOpenTelemetry(
				telemetry.ApplicationName(applicationName),
			)
		}),
		OtelModule,
		LogModule,
		EchoModule,
	)
}

var LogModule = fx.Module("telemetryfx:logger",
	fx.Invoke(telemetry.SetLogger),
)

var OtelModule = fx.Module("telemetryfx:otel",
	fx.Invoke(func(lifecycle fx.Lifecycle, openTelemetry telemetry.OpenTelemetry) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return telemetry.Initialize(ctx, openTelemetry)
			},
			OnStop: telemetry.Shutdown,
		})
	}))

var EchoModule = fx.Module("telemetryfx:echo",
	fx.Provide(func(openTelemetry telemetry.OpenTelemetry) telemetry.EchoMiddlewares {
		return telemetry.NewEchoMiddleware(
			telemetry.Otel(openTelemetry),
			telemetry.SkipHealthCheckRoutes(),
		)
	}))
