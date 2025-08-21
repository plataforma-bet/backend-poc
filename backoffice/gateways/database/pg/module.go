package pg

import (
	"backend-poc/backoffice/extensions/pg"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("pg:repository",
	fx.Provide(pg.New),
	fx.Provide(NewRepository),
	fx.Invoke(func(lc fx.Lifecycle, pool *pg.Pool) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return pool.Close()
			},
		})
	}),
)
