package typed

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Context[T any] struct {
	echo.Context
	ParsedBody T
}

func (r *Context[T]) AddHeader(name string, value string) {
	r.Response().Header().Set(name, value)
}

func (r *Context[T]) Created(location string) error {
	r.AddHeader("Location", location)
	return r.NoContent(http.StatusCreated)
}

func (r *Context[T]) WithStatus(status int) {
	r.Response().WriteHeader(status)
}
