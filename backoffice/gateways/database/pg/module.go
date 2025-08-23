package pg

import (
	"backend-poc/backoffice/extensions/pg"
	"backend-poc/backoffice/gateways/database/pg/session"
	"backend-poc/backoffice/gateways/database/pg/user"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("pg:repository",
	fx.Provide(pg.New),
	fx.Provide(user.NewRepository),
	fx.Provide(session.NewRepository),
	fx.Invoke(func(lc fx.Lifecycle, pool *pg.Pool) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		})
	}),
)
