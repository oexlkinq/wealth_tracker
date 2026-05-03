package tractsgen

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/teambition/rrule-go"
)

type Repo interface {
	ListRtxnsEnds(ctx context.Context) ([]db_api.ListRtxnsEndsRow, error)
	CreateTxn(ctx context.Context, arg db_api.CreateTxnParams) (int64, error)
}

type tractsgen struct {
	repo Repo
}

func New(repos Repo) *tractsgen {
	return &tractsgen{repos}
}

func (v *tractsgen) GenUpTo(ctx context.Context, until time.Time) error {
	rows, err := v.repo.ListRtxnsEnds(ctx)
	if err != nil {
		return err
	}

	for _, row := range rows {
		rr, err := rrule.StrToRRule(row.Rrule)
		if err != nil {
			panic(err)
		}

		if row.Ts.Valid {
			rr.DTStart(row.Ts.Time)
		}

		rr.Until(until)

		for _, occ := range rr.All() {
			_, err := v.repo.CreateTxn(ctx, db_api.CreateTxnParams{
				Ts:     pgtype.Timestamptz{Time: occ},
				Amount: row.Amount,
				RtxnID: pgtype.Int4{Int32: row.ID},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
