package pg

import (
	"backend-poc/backoffice/config"
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Pool struct {
	*gorm.DB
}

func New(ctx context.Context, conStr config.PostgresSQL) (*Pool, error) {
	return NewFromConnString(ctx, conStr.PoolConnectionString())
}

func NewFromConnString(ctx context.Context, connStr string) (*Pool, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Pool{DB: db}, nil
}

func (p *Pool) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
