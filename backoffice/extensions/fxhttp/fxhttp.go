package fxhttp

import (
	"backend-poc/backoffice/config"
	"backend-poc/backoffice/extensions/telemetry"
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:fxhttp",
	fx.Provide(echo.New),
	fx.Invoke(func(e *echo.Echo) {
		e.GET("liveness", func(c echo.Context) error {
			return c.String(200, "OK")
		})

		e.GET("readiness", func(c echo.Context) error {
			return c.String(200, "OK")
		})
	}),
	fx.Invoke(func(echoMiddleware telemetry.EchoMiddlewares, e *echo.Echo) {
		e.Use(echoMiddleware.All()...)
		e.Use(middleware.CORS(), middleware.Recover())
	}),
	fx.Invoke(StartServer),
)

func StartServer(lc fx.Lifecycle, e *echo.Echo, config config.Global) *echo.Echo {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				slog.Info("starting server", "address", config.Address)

				err := e.Start(config.Address)
				if err != nil {
					slog.Error("error starting server", "error", err)
					return
				}

				slog.Info("server started")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
