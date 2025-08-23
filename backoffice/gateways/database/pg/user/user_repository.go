package user

import (
	"backend-poc/backoffice/domain/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	queryCreteUser = `
		INSERT INTO users (id, name, email, password, role, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	queryGetUserByEmail = `
    	SELECT id, name, email, password, role, active, created_at, updated_at
		FROM users
    	WHERE email = $1
	`
)

func (r *Repository) Create(ctx context.Context, user *user.User) error {
	const operation = "user.Repository.CreateUser"

	_, err := r.Exec(
		ctx,
		queryCreteUser,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Active,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const operation = "user.Repository.GetUserByEmail"

	var result user.User

	err := r.QueryRow(ctx, queryGetUserByEmail, email).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Password,
		&result.Role,
		&result.Active,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", operation, user.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &result, nil
}
