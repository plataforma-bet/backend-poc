package config

import "fmt"

type Global struct {
	APIName string `conf:"default:backoffice:api"`
	Address string `conf:"env:SERVER_ADDR,default::3000"`
}

type PostgresSQL struct {
	DatabaseName     string `conf:"env:NAME,default:postgres"`
	DatabaseUser     string `conf:"env:APP_USER,default:postgres"`
	DatabasePassword string `conf:"env:APP_PASSWORD,default:postgres,mask"`
	DatabaseHost     string `conf:"env:ENDPOINT,default:localhost"`
	DatabasePort     string `conf:"env:PORT,default:5432"`
	DatabaseSSLMode  string `conf:"env:SSLMODE,default:disable"`

	DatabaseMigrationUser     string `conf:"env:MIGRATION_USER,default:postgres"`
	DatabaseMigrationPassword string `conf:"env:MIGRATION_PASSWORD,default:postgres,mask"`
}

type ConfigJWT struct {
	Secret        string `conf:"env:SECRET_JWT,default:secret"`
	AccessExpiry  int    `conf:"env:ACCESS_EXPIRY,default:15"`
	RefreshExpiry int    `conf:"env:REFRESH_EXPIRY,default:7"`
}

func (c *PostgresSQL) PoolConnectionString() string {
	if c.DatabaseSSLMode == "" {
		c.DatabaseSSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort,
		c.DatabaseName, c.DatabaseSSLMode)
}

func (c *PostgresSQL) MigrationConnectionString() string {
	if c.DatabaseSSLMode == "" {
		c.DatabaseSSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DatabaseMigrationUser, c.DatabaseMigrationPassword, c.DatabaseHost, c.DatabasePort,
		c.DatabaseName, c.DatabaseSSLMode)
}
