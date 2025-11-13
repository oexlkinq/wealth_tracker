package main

import (
	"database/sql"
	"log"

	"github.com/oexlkinq/wealth_tracker/internal/config"
	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	_ "modernc.org/sqlite"
)

func run() error {
	cfg := config.Load()

	db, err := sql.Open("sqlite", cfg.DB_FILE)
	if err != nil {
		return err
	}

	db_api.New(db)
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}
