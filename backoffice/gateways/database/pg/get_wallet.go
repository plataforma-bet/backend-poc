package pg

import (
	"backend-poc/backoffice/domain/wallet"
	"context"
	"errors"
	"fmt"
)

func (r Repository) GetWalletByUserId(ctx context.Context, userId uint) (*wallet.Wallet, error) {
	const operation = "pg.Repository.GetWalletByUserId"

	var userWallet wallet.Wallet

	if err := r.WithContext(ctx).First(&userWallet, userId).Error; err != nil {
		if errors.Is(err, r.ErrNotFound()) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &userWallet, nil
}
