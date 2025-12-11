package tractsiter

import (
	"context"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/dbiter"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/gen"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/models"
)

type TractsIter struct {
	*dbiter.DBIter
	*gen.GTractGen
}

func New(ctx context.Context, qtx *db_api.Queries, since time.Time, rtract *db_api.Rtract) (*TractsIter, error) {
	i, err := dbiter.New(ctx, qtx, since, rtract)
	if err != nil {
		return nil, err
	}

	g, err := gen.New(rtract, since)
	if err != nil {
		return nil, err
	}

	return &TractsIter{
		DBIter:    i,
		GTractGen: g,
	}, nil
}

func (v *TractsIter) All() iter.Seq[*models.CalcTract] {
	return func(yield func(*models.CalcTract) bool) {
		var lastDate time.Time
		next, stop := iter.Pull(v.DBIter.All())
		for {
			tract, ok := next()
			if !ok {
				break
			}
			lastDate = tract.Date

			if !yield(tract) {
				stop()
				return
			}
		}

		if !lastDate.IsZero() {
			v.GTractGen.ResetTo(lastDate, false)
		}

		for tract := range v.GTractGen.All() {
			if !yield(tract) {
				return
			}
		}
	}
}
