package telemetry

import (
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const scope = "github.com/plataforma-bet/backend-poc/backoffice/extensions/telemetry"

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

type MiddlewaresOption func(*EchoMiddlewares)

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

func (m EchoMiddlewares) Logger() echo.MiddlewareFunc {
	loggerCfg := middleware.DefaultLoggerConfig
	loggerCfg.Skipper = m.skipper

	mi := middleware.LoggerWithConfig(loggerCfg)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return mi(next)(c)
		}
	}
}

func (m EchoMiddlewares) Duration() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	recorder, err := meter.Float64Histogram(
		"http.server.duration",
		metric.WithExplicitBucketBoundaries(0, 0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
	)
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			if err := next(c); err != nil {
				return err
			}

			recorder.Record(c.Request().Context(),
				time.Since(start).Seconds(),
				metric.WithAttributes(attribute.String("http.route", c.Path())),
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
				metric.WithAttributes(attribute.Int("http.response.status_code", c.Response().Status)),
			)

			return nil
		}
	}
}

func (m EchoMiddlewares) RequestCounter() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	recorder, err := meter.Int64Counter("http.server.request.count")
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				return err
			}

			recorder.Add(c.Request().Context(),
				1,
				metric.WithAttributes(attribute.String("http.route", c.Path())),
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
				metric.WithAttributes(attribute.Int("http.response.status_code", c.Response().Status)))

			return nil
		}
	}
}

func (m EchoMiddlewares) ActiveRequests() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	activeRequests, err := meter.Int64UpDownCounter("http.server.active_request")
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			activeRequests.Add(c.Request().Context(),
				1,
				metric.WithAttributes(attribute.String("http.route", c.Path())),
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
			)

			defer activeRequests.Add(c.Request().Context(),
				-1,
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
				metric.WithAttributes(attribute.Int("http.response.status_code", c.Response().Status)),
			)

			return next(c)
		}
	}
}

func (m EchoMiddlewares) RequestSizeHistogram() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	recorder, err := meter.Int64Histogram("http.server.request.size")
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqSize := c.Request().Header.Get("content-length")
			if reqSize != "" {
				reqSize = "0"
			}

			size, _ := strconv.ParseInt(reqSize, 10, 64)

			recorder.Record(c.Request().Context(),
				size,
				metric.WithAttributes(attribute.String("http.route", c.Path())),
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
			)

			return next(c)
		}
	}
}

func (m EchoMiddlewares) ResponseSizeHistogram() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	recorder, err := meter.Int64Histogram("http.server.response.size")
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				return err
			}

			respSize := c.Response().Size

			recorder.Record(c.Request().Context(),
				respSize,
				metric.WithAttributes(attribute.String("http.route", c.Path())),
				metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
				metric.WithAttributes(attribute.Int("http.response.status_code", c.Response().Status)),
			)

			return nil
		}
	}
}

func (m EchoMiddlewares) ErrorCounter() echo.MiddlewareFunc {
	meter := otel.Meter(scope)

	recorder, err := meter.Int64Counter("http.server.request.errors")
	if err != nil {
		return failWith(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				recorder.Add(c.Request().Context(), 1,
					metric.WithAttributes(attribute.String("http.route", c.Path())),
					metric.WithAttributes(attribute.String("http.request.method", c.Request().Method)),
					metric.WithAttributes(attribute.String("error.type", reflect.TypeOf(err).String())),
					metric.WithAttributes(attribute.Int("http.response.status_code", c.Response().Status)),
				)

				return err
			}

			return nil
		}
	}
}

func failWith(err error) func(_ echo.HandlerFunc) echo.HandlerFunc {
	return func(_ echo.HandlerFunc) echo.HandlerFunc {
		return func(_ echo.Context) error {
			return err
		}
	}
}

func (m EchoMiddlewares) All() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		m.RequestID(),
		m.Tracer(),
	}
}
