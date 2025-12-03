package app

import (
	"context"
	"database/sql"

	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

type App struct {
	DB      *sql.DB
	Queries *db_api.Queries
	Tx      *sql.Tx
}

func New(ctx context.Context) (*App, error) {
	// TODO: вынести это в .env или в конфиг и юзать viper
	DBPath := "wealth_tracker.db"

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, err
	}

	queries := db_api.New(db)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	qtx := queries.WithTx(tx)

	return &App{db, qtx, tx}, nil
}

type runEFunc func(cmd *cobra.Command, args []string) error

func MakeCmdRunEFunc(handler func(cmd *cobra.Command, args []string, ctx context.Context, app *App) error) runEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := New(ctx)
		if err != nil {
			return err
		}

		err = handler(cmd, args, ctx, app)
		if err != nil {
			return err
		}

		err = app.Tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}
}
