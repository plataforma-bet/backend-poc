package wallet

import "errors"

var ErrUserWalletNotFound = errors.New("user wallet not found")

var InsufficientBalance = errors.New("insufficient balance")
