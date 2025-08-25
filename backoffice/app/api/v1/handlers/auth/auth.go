package auth

import (
	"backend-poc/backoffice/app/api/v1/request"
	"backend-poc/backoffice/app/api/v1/typed"
	"backend-poc/backoffice/extensions/telemetry"
	"backend-poc/backoffice/usecases/user"
	"net/http"
)

func (h *Handler) Register(ctx typed.Context[request.RegisterRequest]) error {
	input := createRegisterInputFromRequest(ctx.ParsedBody)

	audit := user.Audit{
		UserAgent: ctx.Request().UserAgent(),
		IPAddress: ctx.RealIP(),
	}

	response, err := h.registerUseCase.Register(ctx.Request().Context(), input, audit)
	if err != nil {
		telemetry.HandlerUnexpectedError(ctx.Request().Context(), err)
		return err
	}

	return ctx.JSON(http.StatusCreated, response)
}
