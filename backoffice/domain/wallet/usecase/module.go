package usecase

import "go.uber.org/fx"

var Module = fx.Module("wallet:usecase",
	fx.Provide(
		NewWallet,
	),
)
