package user

import (
	"backend-poc/backoffice/domain/session"
	"backend-poc/backoffice/domain/user"
	"backend-poc/backoffice/extensions/auth"
	"backend-poc/backoffice/extensions/telemetry"
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *UseCase) Register(ctx context.Context, input RegisterInput, audit Audit) (*RegisterOutPut, error) {
	const operation = "UseCase.Auth.Register"

	ctx, span := telemetry.StartSpan(ctx, operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("id", input.ID.String()),
		attribute.String("role", input.Role),
		attribute.String("user_agent", audit.UserAgent),
		attribute.String("ip_address", audit.IPAddress),
	)

	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	foundUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if !errors.Is(err, user.ErrUserNotFound) && foundUser != nil {
		return nil, fmt.Errorf("%s: %w", operation, ErrUserAlreadyExists)
	}

	newUser := user.NewUser(input.Name, input.Email, input.Password, input.Role)

	if err = newUser.HashPassword(); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if err = uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	token, err := uc.auth.GenerateAccessToken(auth.Subject{
		UserID: newUser.ID,
		Email:  newUser.Email,
		Role:   newUser.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	refreshToken, err := uc.auth.GenerateRefreshToken(auth.Subject{
		UserID: newUser.ID,
		Email:  newUser.Email,
		Role:   newUser.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	newSession := session.NewSession(
		newUser.ID,
		token,
		refreshToken,
		time.Now().Add(uc.auth.GetAccessTokenExpiry()),
		audit.IPAddress,
		audit.UserAgent,
	)

	if err = uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &RegisterOutPut{
		User: UserResponse{
			ID:        newUser.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			Role:      newUser.Role,
			Active:    newUser.Active,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		},
		AccessToken:  token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(uc.auth.GetAccessTokenExpiry()),
		TokenType:    token,
	}, nil
}
