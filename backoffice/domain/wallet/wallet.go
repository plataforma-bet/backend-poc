package wallet

import (
	"backend-poc/backoffice/domain"
	"backend-poc/backoffice/domain/auth"
	"time"
)

type Wallet struct {
	domain.Entity
	UserId         uint      `json:"user_id" gorm:"type:integer;column:user_id;not_null;index"`
	Balance        int       `json:"balance" gorm:"not_null;default:0;index"`
	BlockedBalance int       `json:"blocked_balance" gorm:"not_null;default:0;index"`
	BonusBalance   int       `json:"bonus_balance" gorm:"not_null;default:0;index"`
	Currency       string    `json:"currency" gorm:"type:varchar(10);not_null;default:'BRL'"`
	Country        string    `json:"country" gorm:"type:varchar(10);not_null;default:BR"`
	User           auth.User `json:"-" gorm:"foreignKey:UserId"`
}

func NewWallet(userId uint, currency *string, country *string) *Wallet {
	return &Wallet{
		Currency:       *currency,
		Country:        *country,
		UserId:         userId,
		BlockedBalance: 0,
		BonusBalance:   0,
		Balance:        0,
	}
}

func (u *Wallet) IncreaseBalance(amount int) {
	u.Balance += amount

	u.UpdatedAt = time.Now()
}
func (u *Wallet) DecreaseBalance(amount int) error {

	if amount > u.Balance {
		return InsufficientBalance
	}

	u.Balance -= amount

	u.UpdatedAt = time.Now()

	return nil
}
