package handlers

import (
	userAuth "backend-poc/backoffice/app/api/v1/handlers/auth"

	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api:handlers",
	fx.Provide(userAuth.NewHandler),
)
