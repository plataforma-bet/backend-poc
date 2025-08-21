package v1

import (
	"backend-poc/backoffice/app/api/v1/handlers"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("backoffice:api:v1",
	handlers.Module,
	fx.Invoke(func(e *echo.Echo, handler handlers.Handler) {
		v1 := e.Group("/api/v1")

		v1.POST("/users", nil)
		v1.GET("/users/:id", nil)
		v1.PUT("/users/:id", nil)
		v1.DELETE("/users/:id", nil)
	}))
