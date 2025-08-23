package main

import (
	"backend-poc/backoffice/config"
	"backend-poc/backoffice/extensions/telemetry/telemetryfx"
	"backend-poc/backoffice/gateways/database/pg/migration"
	"context"
	"log/slog"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(context.Background),
		config.Module,
		telemetryfx.Module("backoffice:migrate"),
		fx.Invoke(func(sql config.PostgresSQL, lifecycle fx.Lifecycle, shutdower fx.Shutdowner) {
			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						err := migration.Migrate(sql.MigrationConnectionString())
						if err != nil {
							slog.Error("error trying to run migrations", slog.Any("error", err))
							_ = shutdower.Shutdown(fx.ExitCode(1))

							return
						}
						_ = shutdower.Shutdown(fx.ExitCode(1))
					}()

					return nil
				},
			})
		}),
	).Run()
}
