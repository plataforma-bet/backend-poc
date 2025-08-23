package modules

import (
	routes "backend-poc/backoffice/app/api/v1"
	"backend-poc/backoffice/config"
	"backend-poc/backoffice/extensions/fxhttp"
	"backend-poc/backoffice/extensions/telemetry/telemetryfx"
	pggate "backend-poc/backoffice/gateways/database/pg"
	"backend-poc/backoffice/usecases"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api",
	fx.Provide(context.Background),
	config.Module,
	pggate.Module,
	telemetryfx.Module("backoffice:api"),
	usecases.Module,
	routes.Module,
	fxhttp.Module,
)
