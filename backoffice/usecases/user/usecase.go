package user

import (
	"backend-poc/backoffice/extensions/auth"
	"backend-poc/backoffice/gateways/database/pg/session"
	"backend-poc/backoffice/gateways/database/pg/user"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("password must be at least 8 characters long")
	ErrInvalidName       = errors.New("name must be at least 2 characters long")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type RegisterInput struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Role     string
	Active   bool
}

type Audit struct {
	UserAgent string
	IPAddress string
}

type RegisterOutPut struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	TokenType    string       `json:"token_type"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UseCase struct {
	userRepo    *user.Repository
	sessionRepo *session.Repository
	auth        *auth.JWTServices
}

func NewUseCase(user *user.Repository, session *session.Repository, auth *auth.JWTServices) *UseCase {
	return &UseCase{
		userRepo:    user,
		sessionRepo: session,
		auth:        auth,
	}
}

func (r *RegisterInput) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return ErrInvalidEmail
	}

	if len(r.Password) < 8 {
		return ErrInvalidPassword
	}

	if len(r.Name) < 2 {
		return ErrInvalidName
	}

	return nil
}
