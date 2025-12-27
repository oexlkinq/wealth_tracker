package calc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/tractsgen"
)

type TargetReachInfo struct {
	Target        *db_api.Target
	ReachDate     time.Time
	ReachedAmount float64
	Reached       bool
}

func Calc(ctx context.Context, qtx *db_api.Queries) ([]*TargetReachInfo, error) {
	targets, err := qtx.ListTargetsForCalc(ctx)
	if err != nil {
		return nil, err
	}

	generatedUntil := time.Now()

	tris := make([]*TargetReachInfo, len(targets))
	for i, target := range targets {
		var date time.Time
		for i := range 100 {
			date, err = qtx.GetReachingTargetDate(ctx, target.Amount)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					if i == 100-1 {
						panic("too many retries")
					}

					generatedUntil = generatedUntil.Add(time.Hour * 24 * 365)

					err = tractsgen.GenUpTo(ctx, qtx, generatedUntil)
					if err != nil {
						return nil, err
					}

					fmt.Println("pushed until", generatedUntil)
					continue
				}

				return nil, err
			}

			break
		}

		tractId, err := qtx.CreateTract(ctx, db_api.CreateTractParams{
			Type:   "target",
			Date:   date,
			Amount: -target.Amount,
			Acked:  false,
		})
		if err != nil {
			return nil, err
		}

		err = qtx.UpdateTractIDOfTarget(ctx, db_api.UpdateTractIDOfTargetParams{
			TractID:  sql.NullInt64{Int64: tractId},
			TargetID: target.ID,
		})
		if err != nil {
			return nil, err
		}

		tris[i] = &TargetReachInfo{
			Target:        &targets[i],
			ReachDate:     date,
			ReachedAmount: target.Amount,
			Reached:       true,
		}
	}

	return tris, nil
}
