package typed

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var sideEffectMethods = []string{
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
}

func Validated[T any](handler func(request Context[T]) error) echo.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return Typed(func(request Context[T]) error {
		if err := validate.Struct(request.ParsedBody); err != nil {
			slog.Debug("body failed validation", slog.Any("error", err))

			request.Response().WriteHeader(http.StatusUnprocessableEntity)
			return err
		}

		return handler(request)
	})
}

func Typed[T any](handler func(request Context[T]) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		var t T

		if slices.Contains(sideEffectMethods, e.Request().Method) {
			if err := e.Bind(&t); err != nil {
				e.Response().WriteHeader(http.StatusBadRequest)
				return err
			}
		}

		return handler(Context[T]{Context: e, ParsedBody: t})
	}
}
