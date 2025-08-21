package modules

import (
	"backend-poc/backoffice/config"
	"backend-poc/backoffice/extensions/fxhttp"
	"backend-poc/backoffice/extensions/telemetry/telemetryfx"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api",
	fx.Provide(context.Background),
	telemetryfx.Module("backoffice:api"),
	fxhttp.Module,
	config.Module,
)
