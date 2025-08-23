package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Token        string    `json:"token" db:"token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	IsRevoked    bool      `json:"is_revoked" db:"is_revoked"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewSession(userID uuid.UUID, token, refreshToken string, expiresAt time.Time, ipAddress, userAgent string) *Session {
	now := time.Now()
	return &Session{
		ID:           uuid.New(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		IsRevoked:    false,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Revoke() {
	s.IsRevoked = true
	s.UpdatedAt = time.Now()
}

func (s *Session) IsValid() bool {
	return !s.IsRevoked && !s.IsExpired()
}
