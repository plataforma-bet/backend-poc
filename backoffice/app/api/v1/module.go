package v1

import (
	"backend-poc/backoffice/app/api/v1/handlers"
	"backend-poc/backoffice/app/api/v1/handlers/auth"
	"backend-poc/backoffice/app/api/v1/request"
	"backend-poc/backoffice/app/api/v1/typed"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api:v1",
	handlers.Module,
	fx.Invoke(func(e *echo.Echo, handler auth.Handler) {
		v1 := e.Group("/api/v1")

		a := v1.Group("/auth")
		a.POST("/register", typed.Validated[request.RegisterRequest](handler.Register))
	}),
)
