package calc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
)

type TargetReachInfo struct {
	Target    db_api.Target
	ReachDate time.Time
}

type Repo interface {
	ListTargetsForCalc(ctx context.Context) ([]db_api.Target, error)
	GetReachingTargetDate(ctx context.Context, amount float64) (time.Time, error)
	CreateTract(ctx context.Context, arg db_api.CreateTractParams) (int64, error)
}

type Tractsgen interface {
	GenUpTo(ctx context.Context, until time.Time) error
}

func Calc(ctx context.Context, repo Repo, tg Tractsgen) ([]TargetReachInfo, error) {
	targets, err := repo.ListTargetsForCalc(ctx)
	if err != nil {
		return nil, err
	}

	generatedUntil := time.Now().Truncate(time.Hour * 24)

	tris := make([]TargetReachInfo, len(targets))
	for i, target := range targets {
		var date time.Time
		for i := range 100 {
			if i == 100-1 {
				panic("too many retries")
			}

			date, err = repo.GetReachingTargetDate(ctx, target.Amount)
			if err != nil {
				return nil, err
			}

			if !date.IsZero() {
				break
			}

			generatedUntil = generatedUntil.Add(time.Hour * 24 * 365)

			err = tg.GenUpTo(ctx, generatedUntil)
			if err != nil {
				return nil, err
			}

			fmt.Println("pushed until", generatedUntil)
		}

		_, err := repo.CreateTract(ctx, db_api.CreateTractParams{
			Date:     date,
			Amount:   target.Amount,
			Acked:    false,
			TargetID: sql.NullInt64{Int64: target.ID, Valid: true},
		})
		if err != nil {
			return nil, err
		}

		tris[i] = TargetReachInfo{
			Target:    targets[i],
			ReachDate: date,
		}
	}

	return tris, nil
}
