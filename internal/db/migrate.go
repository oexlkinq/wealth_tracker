package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/oexlkinq/wealth_tracker/internal/config"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate() error {
	db, err := sql.Open("pgx", config.POSTGRES_DBSTRING)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
