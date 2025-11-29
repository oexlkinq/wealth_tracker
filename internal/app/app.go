package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	_ "modernc.org/sqlite"
)

type App struct {
	DB      *sql.DB
	Queries *db_api.Queries
}

func New() (*App, error) {
	DBPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		return nil, fmt.Errorf("env var DB_PATH is not defined")
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, err
	}

	queries := db_api.New(db)

	return &App{db, queries}, nil
}
