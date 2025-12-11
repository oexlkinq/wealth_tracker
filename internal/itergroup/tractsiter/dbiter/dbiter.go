package dbiter

import (
	"context"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/models"
)

type DBIter struct {
	calcTracts []*models.CalcTract
}

func New(ctx context.Context, qtx *db_api.Queries, since time.Time, rtract *db_api.Rtract) (*DBIter, error) {
	tracts, err := qtx.ListTractsSince(ctx, db_api.ListTractsSinceParams{
		Since:    since,
		RtractID: rtract.ID,
	})
	if err != nil {
		return nil, err
	}

	calcTracts := make([]*models.CalcTract, len(tracts))
	for i := range tracts {
		calcTracts[i] = &models.CalcTract{
			Tract:     &tracts[i],
			RTractID:  0,
			Generated: false,
		}
	}

	return &DBIter{calcTracts: calcTracts}, nil
}

func (v *DBIter) All() iter.Seq[*models.CalcTract] {
	return func(yield func(*models.CalcTract) bool) {
		for i := range v.calcTracts {
			if !yield(v.calcTracts[i]) {
				return
			}
		}
	}
}
