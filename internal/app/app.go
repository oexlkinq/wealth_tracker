package app

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/oexlkinq/wealth_tracker/internal/db"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/spf13/cobra"
)

type App struct {
	DB      *pgx.Conn
	Queries *db_api.Queries
	Tx      pgx.Tx
}

func New(ctx context.Context) (*App, error) {
	// TODO: вынести это в .env или в конфиг и юзать viper
	pgConnStr := os.Getenv("POSTGRES_DBSTRING")

	pgc, err := pgx.Connect(ctx, pgConnStr)
	if err != nil {
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		return nil, err
	}

	queries := db_api.New(pgc)

	tx, err := pgc.Begin(ctx)
	if err != nil {
		return nil, err
	}

	qtx := queries.WithTx(tx)

	return &App{
		DB:      pgc,
		Queries: qtx,
		Tx:      tx,
	}, nil
}

type runEFunc func(cmd *cobra.Command, args []string) error

func MakeCmdRunEFunc(handler func(cmd *cobra.Command, args []string, ctx context.Context, app *App) error) runEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := New(ctx)
		if err != nil {
			return fmt.Errorf("create app: %w", err)
		}
		defer app.Tx.Rollback(ctx)

		err = handler(cmd, args, ctx, app)
		if err != nil {
			return fmt.Errorf("run handler: %w", err)
		}

		err = app.Tx.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}
