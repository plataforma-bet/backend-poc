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
}

func (c *PostgresSQL) PoolConnectionString() string {
	if c.DatabaseSSLMode == "" {
		c.DatabaseSSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort,
		c.DatabaseName, c.DatabaseSSLMode)
}
