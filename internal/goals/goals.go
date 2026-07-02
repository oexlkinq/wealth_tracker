package goals

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
	ListGoalsForCalc(ctx context.Context) ([]db_api.Goal, error)
	GetGoalReachTs(ctx context.Context, target float64) ([]pgtype.Timestamptz, error)
	CreateTxn(ctx context.Context, arg db_api.CreateTxnParams) (int32, error)
}

var _ Repo = (*db_api.Queries)(nil)

type TxnsGen interface {
	GenNextRange(ctx context.Context) error
}

var _ TxnsGen = (*txnsgen.TxnsGen)(nil)

type GoalsSvc struct {
	repo Repo
	tg   TxnsGen
}

func New(repo Repo, tg TxnsGen) (*GoalsSvc, error) {
	return &GoalsSvc{
		repo: repo,
		tg:   tg,
	}, nil
}

func (svc *GoalsSvc) CalcGoals(ctx context.Context) ([]GoalReachInfo, error) {
	goals, err := svc.repo.ListGoalsForCalc(ctx)
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

			tsRows, err := svc.repo.GetGoalReachTs(ctx, goal.Amount)
			if err != nil {
				return nil, err
			}

			// дата достижения цели найдена
			if len(tsRows) != 0 {
				ts = tsRows[0]
				break
			}

			// иначе сгенерить следующую пачку транзакций
			err = svc.tg.GenNextRange(ctx)
			if err != nil {
				return nil, err
			}
		}

		_, err := svc.repo.CreateTxn(ctx, db_api.CreateTxnParams{
			Amount:  -goal.Amount,
			Comment: pgtype.Text{String: "goal done", Valid: true},
			Ts:      ts,
			GoalID:  pgtype.Int4{Int32: goal.ID, Valid: true},
		})
		if err != nil {
			return nil, err
		}

		gris[i] = GoalReachInfo{
			Goal:      goal,
			ReachDate: ts.Time,
		}
	}

	return gris, nil
}
