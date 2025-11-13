package dbiter

import (
	"context"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db_api"
)

type DBIter struct {
	tracts []db_api.ListTractsSinceRow
}

func New(ctx context.Context, qtx *db_api.Queries, since time.Time, rtract *db_api.Rtract) (*DBIter, error) {
	tracts, err := qtx.ListTractsSince(ctx, db_api.ListTractsSinceParams{
		Since:    since,
		RtractID: rtract.ID,
	})
	if err != nil {
		return nil, err
	}

	return &DBIter{tracts}, nil
}

func (v *DBIter) All() iter.Seq[*db_api.ListTractsSinceRow] {
	return func(yield func(*db_api.ListTractsSinceRow) bool) {
		for i := range v.tracts {
			if !yield(&v.tracts[i]) {
				return
			}
		}
	}
}
