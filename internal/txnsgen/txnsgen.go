package txnsgen

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/teambition/rrule-go"
)

type Repo interface {
	ListRtxnsEnds(ctx context.Context) ([]db_api.ListRtxnsEndsRow, error)
	CreateTxn(ctx context.Context, arg db_api.CreateTxnParams) (int32, error)
}

var _ Repo = (*db_api.Queries)(nil)

type TxnsGen struct {
	repo Repo
}

func New(repos Repo) *TxnsGen {
	return &TxnsGen{repos}
}

// default: 2 *basic* years
var RangeStepSize = time.Duration(time.Hour * 24 * 365 * 2)

func (v *TxnsGen) GenNextRange(ctx context.Context) error {
	rows, err := v.repo.ListRtxnsEnds(ctx)
	if err != nil {
		return err
	}

	for _, row := range rows {
		rr, err := rrule.StrToRRule(row.Rrule)
		if err != nil {
			return fmt.Errorf("str to rrule: %w", err)
		}

		// если уже есть сгенеренные транзы, то сместить до последней.
		// иначе будет использован родной dtstart, который обязательно есть в rrule по стандарту
		if row.Ts.Valid {
			rr.DTStart(row.Ts.Time)
		}

		rr.Until(rr.GetDTStart().Add(RangeStepSize))

		for _, nextTime := range rr.All() {
			_, err := v.repo.CreateTxn(ctx, db_api.CreateTxnParams{
				Ts:     pgtype.Timestamptz{Time: nextTime},
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
