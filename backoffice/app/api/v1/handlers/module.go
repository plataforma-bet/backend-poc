package handlers

import "go.uber.org/fx"

var Module = fx.Module("api:v1:handlers",
	fx.Provide(NewHandler),
)
