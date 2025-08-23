package pg

import (
	"backend-poc/backoffice/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func New(ctx context.Context, conStr config.PostgresSQL) (*Pool, error) {
	return NewFromConnString(ctx, conStr.PoolConnectionString())
}

func NewFromConnString(ctx context.Context, connStr string) (*Pool, error) {
	c, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Pool: pool,
	}, nil
}

func (p Pool) Close() {
	p.Pool.Close()
}
