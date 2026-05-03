package calc

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
)

type GoalReachInfo struct {
	Goal      db_api.Goal
	ReachDate time.Time
}

type Repo interface {
	ListGoalsForCalc(ctx context.Context) ([]db_api.Goal, error)
	GetGoalReachTs(ctx context.Context, target float64) ([]pgtype.Timestamptz, error)
	CreateTxn(ctx context.Context, arg db_api.CreateTxnParams) (int64, error)
}

type TxnsGen interface {
	GenUpTo(ctx context.Context, until time.Time) error
}

func Calc(ctx context.Context, repo Repo, tg TxnsGen) ([]GoalReachInfo, error) {
	goals, err := repo.ListGoalsForCalc(ctx)
	if err != nil {
		return nil, err
	}

	generatedUntil := time.Now().Truncate(time.Hour * 24)

	gris := make([]GoalReachInfo, len(goals))
	for i, goal := range goals {
		var ts []pgtype.Timestamptz
		for i := range 100 {
			if i == 100-1 {
				panic("too many retries")
			}

			ts, err = repo.GetGoalReachTs(ctx, goal.Amount)
			if err != nil {
				return nil, err
			}

			if len(ts) == 0 {
				break
			}

			generatedUntil = generatedUntil.Add(time.Hour * 24 * 365)

			err = tg.GenUpTo(ctx, generatedUntil)
			if err != nil {
				return nil, err
			}

			fmt.Println("pushed until", generatedUntil)
		}

		_, err := repo.CreateTxn(ctx, db_api.CreateTxnParams{
			Amount:  goal.Amount,
			Comment: pgtype.Text{String: "TODO", Valid: true},
			Ts:      ts[0],
			RtxnID:  pgtype.Int4{Int32: goal.ID},
		})
		if err != nil {
			return nil, err
		}

		gris[i] = GoalReachInfo{
			Goal:      goals[i],
			ReachDate: ts[0].Time,
		}
	}

	return gris, nil
}
