package itergroup

import (
	"context"
	"iter"
	"slices"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/models"
)

type tractsIterState struct {
	temp *models.CalcTract
	ok   bool
	next func() (*models.CalcTract, bool)
	stop func()
}

type TractsIterGroup []*tractsIterState

func New(ctx context.Context, qtx *db_api.Queries, since time.Time) (TractsIterGroup, error) {
	rtracts, err := qtx.ListRTracts(ctx)
	if err != nil {
		return nil, err
	}

	tig := make(TractsIterGroup, 0, len(rtracts))
	for i := range rtracts {
		ti, err := tractsiter.New(ctx, qtx, since, &rtracts[i])
		if err != nil {
			return nil, err
		}

		next, stop := iter.Pull(ti.All())
		temp, ok := next()

		tig = append(tig, &tractsIterState{
			temp: temp,
			next: next,
			stop: stop,
			ok:   ok,
		})
	}

	return tig, nil
}

func (v TractsIterGroup) All() iter.Seq[*models.CalcTract] {
	return func(yield func(*models.CalcTract) bool) {
		defer v.stop()

		for {
			// отсортировать, чтобы первой оказался итератор с наименьшей датой
			slices.SortFunc(v, func(a, b *tractsIterState) int {
				if !a.ok {
					return 1
				}

				if !b.ok {
					return -1
				}

				return a.temp.Date.Compare(b.temp.Date)
			})
			iter := v[0]

			// если ближайший кончился, то кончились все
			if !iter.ok {
				return
			}

			if !yield(iter.temp) {
				return
			}

			// перевести итератор на следующее значение
			iter.temp, iter.ok = iter.next()
		}
	}
}

func (v TractsIterGroup) stop() {
	for i := range v {
		v[i].stop()
	}
}
