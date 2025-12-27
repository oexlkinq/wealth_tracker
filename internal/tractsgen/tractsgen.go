package tractsgen

import (
	"context"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/teambition/rrule-go"
)

func GenUpTo(ctx context.Context, qtx *db_api.Queries, until time.Time) error {
	rows, err := qtx.ListRtractsWithLastTracts(ctx)
	if err != nil {
		return err
	}

	for _, row := range rows {
		rr, err := rrule.StrToRRule(row.Rrule)
		if err != nil {
			panic(err)
		}

		if row.Date.Valid {
			rr.DTStart(row.Date.Time)
		}

		rr.Until(until)

		for _, occ := range rr.All() {
			tractId, err := qtx.CreateTract(ctx, db_api.CreateTractParams{
				Type:   "rtract",
				Date:   occ,
				Amount: row.Amount,
				Acked:  false,
			})
			if err != nil {
				return err
			}

			err = qtx.CreateRTractToTract(ctx, db_api.CreateRTractToTractParams{
				RtractID: row.ID,
				TractID:  tractId,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
