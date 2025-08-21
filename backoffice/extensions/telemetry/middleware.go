package telemetry

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

type EchoMiddlewares struct {
	openTelemetry OpenTelemetry
	skipper       middleware.Skipper
}

func NewEchoMiddleware(opts ...MiddlewaresOption) EchoMiddlewares {
	m := &EchoMiddlewares{}

	for _, opt := range opts {
		opt(m)
	}

	return *m
}

type MiddlewaresOption func(middlewares *EchoMiddlewares)

func Otel(openTelemetry OpenTelemetry) MiddlewaresOption {
	return func(c *EchoMiddlewares) {
		c.openTelemetry = openTelemetry
	}
}

func Skipper(sk middleware.Skipper) MiddlewaresOption {
	return func(c *EchoMiddlewares) {
		c.skipper = sk
	}
}

func SkipHealthCheckRoutes() MiddlewaresOption {
	return Skipper(func(c echo.Context) bool {
		switch c.Path() {
		case "/liveness", "/readiness":
			return true
		default:
			return false
		}
	})
}

func (m EchoMiddlewares) Tracer() echo.MiddlewareFunc {
	return otelecho.Middleware(
		m.openTelemetry.ApplicationName,
		otelecho.WithTracerProvider(otel.GetTracerProvider()),
		otelecho.WithSkipper(m.skipper),
	)
}

func (m EchoMiddlewares) RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get("x-request-id")
			ctx := ContextWithRequestID(c.Request().Context(), reqID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (m EchoMiddlewares) All() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		m.RequestID(),
		m.Tracer(),
	}
}
