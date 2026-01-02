package tractsgen

import (
	"context"
	"database/sql"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/teambition/rrule-go"
)

type Repo interface {
	ListRtractsWithLastTracts(ctx context.Context) ([]db_api.ListRtractsWithLastTractsRow, error)
	CreateTract(ctx context.Context, arg db_api.CreateTractParams) (int64, error)
}

type tractsgen struct {
	repo Repo
}

func New(repos Repo) *tractsgen {
	return &tractsgen{repos}
}

func (v *tractsgen) GenUpTo(ctx context.Context, until time.Time) error {
	rows, err := v.repo.ListRtractsWithLastTracts(ctx)
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
			_, err := v.repo.CreateTract(ctx, db_api.CreateTractParams{
				Date:     occ,
				Amount:   row.Amount,
				Acked:    false,
				RtractID: sql.NullInt64{Int64: row.ID, Valid: true},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
