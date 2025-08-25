package session

import (
	"backend-poc/backoffice/domain/session"
	"backend-poc/backoffice/extensions/telemetry"
	"context"
	"fmt"
)

const (
	queryCreateSession = `
		INSERT INTO sessions (id, user_id, token, refresh_token, expires_at, is_revoked, ip_address, user_agent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
)

func (r *Repository) Create(ctx context.Context, session *session.Session) error {
	const operation = "session.Repository.Create"

	ctx, span := telemetry.StartSpan(ctx, operation)
	defer span.End()

	_, err := r.Exec(
		ctx,
		queryCreateSession,
		session.ID,
		session.UserID,
		session.Token,
		session.RefreshToken,
		session.ExpiresAt,
		session.IsRevoked,
		session.IPAddress,
		session.UserAgent,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
