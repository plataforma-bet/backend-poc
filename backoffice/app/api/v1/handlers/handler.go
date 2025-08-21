package handlers

import "go.uber.org/fx"

type HandlerParams struct {
	fx.In
}
type Handler struct {
}

func NewHandler(params HandlerParams) Handler {
	return Handler{}
}
