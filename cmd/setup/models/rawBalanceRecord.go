package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
)

type RawBalanceRecord struct {
	Amount float64
	Date   string
}

func (v *RawBalanceRecord) Insert(ctx context.Context, qtx *db_api.Queries) error {
	date, err := time.Parse(time.DateOnly, v.Date)
	if err != nil {
		return err
	}

	err = qtx.CreateBalanceRecord(ctx, db_api.CreateBalanceRecordParams{
		Amount:      v.Amount,
		Date:        date,
		OriginTract: sql.NullInt64{},
	})
	if err != nil {
		return err
	}

	return nil
}
