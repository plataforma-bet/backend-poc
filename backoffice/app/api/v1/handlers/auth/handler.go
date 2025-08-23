package auth

import (
	"backend-poc/backoffice/app/api/v1/request"
	"backend-poc/backoffice/usecases/user"
)

type Handler struct {
	registerUseCase *user.UseCase
}

func NewHandler(register *user.UseCase) Handler {
	return Handler{
		registerUseCase: register,
	}
}

func createRegisterInputFromRequest(r request.RegisterRequest) user.RegisterInput {
	return user.RegisterInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Role:     r.Role,
	}
}
