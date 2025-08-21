package config

import "fmt"

type Global struct {
	APIName string `conf:"default:backoffice:api"`
	Address string `conf:"env:SERVER_ADDR,default::3000"`
}

type PostgresSQL struct {
	DatabaseName        string `conf:"env:NAME,default:db"`
	DatabaseUser        string `conf:"env:APP_USER,default:postgres"`
	DatabasePassword    string `conf:"env:APP_PASSWORD,default:postgres,mask"`
	DatabaseHost        string `conf:"env:ENDPOINT,default:localhost"`
	DatabasePort        string `conf:"env:PORT,default:5432"`
	DatabaseSSLMode     string `conf:"env:SSLMODE,default:disable"`
	DatabasePoolMinSize int32  `conf:"env:POOL_MIN_SIZE,default:2"`
	DatabasePoolMaxSize int32  `conf:"env:POOL_MAX_SIZE,default:10"`
}

func (c *PostgresSQL) PoolConnectionString() string {
	if c.DatabaseSSLMode == "" {
		c.DatabaseSSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_min_conns=%d&pool_max_conns=%d",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort,
		c.DatabaseName, c.DatabaseSSLMode, c.DatabasePoolMinSize, c.DatabasePoolMaxSize)
}
