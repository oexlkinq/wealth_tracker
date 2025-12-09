package tractsiter

import (
	"context"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/dbiter"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/gen"
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

type Tract struct {
	db_api.ListTractsSinceRow
	Generated bool
}

func (v *TractsIter) All() iter.Seq[*Tract] {
	return func(yield func(*Tract) bool) {
		for t := range v.DBIter.All() {
			if !yield(&Tract{
				ListTractsSinceRow: *t,
				Generated:          false,
			}) {
				return
			}
		}

		for tract := range v.GTractGen.All() {
			if !yield(&Tract{
				ListTractsSinceRow: *tract,
				Generated:          true,
			}) {
				return
			}
		}
	}
}
