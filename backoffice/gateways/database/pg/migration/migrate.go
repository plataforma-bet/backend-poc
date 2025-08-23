package migration

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func Migrate(conn string) error {
	d, err := iofs.New(fs, ".")
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, conn)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			slog.Error("closing the source", slog.Any("error", sourceErr))
		}

		if dbErr != nil {
			slog.Error("closing postgres connection", slog.Any("error", dbErr))
		}
	}()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running up migrations: %w", err)
	}

	return nil

}
