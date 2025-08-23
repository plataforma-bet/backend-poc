package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role" `
	Active    bool       `json:"active"`
	LastLogin *time.Time `json:"last_login,omitempty" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewUser(name, email, password, role string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
