package usecases

import (
	"backend-poc/backoffice/extensions/auth"
	"backend-poc/backoffice/usecases/user"

	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:usecases",
	fx.Provide(user.NewUseCase),
	auth.Module,
)
