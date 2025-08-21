package handlers

import (
	"backend-poc/backoffice/domain/wallet/usecase"

	"go.uber.org/fx"
)

type HandlerParams struct {
	fx.In

	GetWalletByUserId *usecase.Wallet
}
type Handler struct {
	getWalletByUserIdUC *usecase.Wallet
}

func NewHandler(params HandlerParams) Handler {
	return Handler{
		getWalletByUserIdUC: params.GetWalletByUserId,
	}
}
