package modules

import (
	v1 "backend-poc/backoffice/app/api/v1"
	"backend-poc/backoffice/config"
	"backend-poc/backoffice/domain"
	"backend-poc/backoffice/extensions/fxhttp"
	"backend-poc/backoffice/extensions/telemetry/telemetryfx"
	pggate "backend-poc/backoffice/gateways/database/pg"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api",
	fx.Provide(context.Background),
	telemetryfx.Module("backoffice:api"),
	domain.Module,
	pggate.Module,
	v1.Module,
	fxhttp.Module,
	config.Module,
)
