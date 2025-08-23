package auth

import "go.uber.org/fx"

var Module = fx.Module("backoffice:extensions:auth",
	fx.Provide(NewJWT),
)
