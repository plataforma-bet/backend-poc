package domain

import (
	"backend-poc/backoffice/domain/wallet/usecase"

	"go.uber.org/fx"
)

var Module = fx.Module("domain:usecases",
	fx.Options(usecase.Module),
)
