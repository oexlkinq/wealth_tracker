package calc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/txnsgen"
)

type GoalReachInfo struct {
	Goal      db_api.Goal
	ReachDate time.Time
}

type Repo interface {
	ListGoalsForCalc(ctx context.Context) ([]db_api.ListGoalsForCalcRow, error)
	GetGoalReachTs(ctx context.Context, target float64) ([]pgtype.Timestamptz, error)
	CreateTxn(ctx context.Context, arg db_api.CreateTxnParams) (int32, error)
}

var _ Repo = (*db_api.Queries)(nil)

type TxnsGen interface {
	GenNextRange(ctx context.Context) error
}

var _ TxnsGen = (*txnsgen.TxnsGen)(nil)

func Calc(ctx context.Context, repo Repo, tg TxnsGen) ([]GoalReachInfo, error) {
	goals, err := repo.ListGoalsForCalc(ctx)
	if err != nil {
		return nil, err
	}

	gris := make([]GoalReachInfo, len(goals))
	for i, goal := range goals {
		var ts pgtype.Timestamptz

		for i := range 100 {
			if i == 100-1 {
				panic("too many retries")
			}

			tsRows, err := repo.GetGoalReachTs(ctx, goal.Amount.Float64)
			if err != nil {
				return nil, err
			}

			// дата достижения цели найдена
			if len(tsRows) != 0 {
				ts = tsRows[0]
				break
			}

			// иначе сгенерить следующую пачку транзакций
			err = tg.GenNextRange(ctx)
			if err != nil {
				return nil, err
			}
		}

		_, err := repo.CreateTxn(ctx, db_api.CreateTxnParams{
			Amount:  goal.Amount.Float64,
			Comment: pgtype.Text{String: "goal done", Valid: true},
			Ts:      ts,
			RtxnID:  goal.ID,
		})
		if err != nil {
			return nil, err
		}

		gris[i] = GoalReachInfo{
			Goal: db_api.Goal{
				ID:      goal.ID.Int32,
				Amount:  goal.Amount.Float64,
				Comment: goal.Comment,
				Index:   goal.Index.Int32,
			},
			ReachDate: ts.Time,
		}
	}

	return gris, nil
}
