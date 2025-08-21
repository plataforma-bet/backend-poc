package config

import "go.uber.org/fx"

var Module = fx.Module("backoffice:config",
	fx.Provide(Config[Global]()),
	fx.Provide(Config[PostgresSQL](Prefix("POSTGRES_DATABASE"))),
)
