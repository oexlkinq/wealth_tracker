package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
)

type DBTX interface {
	db_api.DBTX
	Get(dest any, query string, args ...any) error
}

type Repo struct {
	*db_api.Queries
	db DBTX
}

func NewRepo(db DBTX) *Repo {
	return &Repo{
		Queries: db_api.New(db),
		db:      db,
	}
}

const getReachingTargetDate = `
with
	balance as (select amount, date from balance_records order by date desc limit 1),
	csums as (
		select id, amount, date, (select b.amount from balance b) + sum(amount) OVER (ORDER BY date, id ROWS UNBOUNDED PRECEDING) as csum
		from tracts
		where date > (select date from balance)
	),
	csums_with_prev as (
		select id, amount, date, csum, lag(csum, 1) over (order by date, id) as prev_csum
		from csums
	)
select date
from csums_with_prev
where csum >= ? and (prev_csum < ? or prev_csum is null)
order by date desc, id desc
limit 1`

func (v *Repo) GetReachingTargetDate(ctx context.Context, amount float64) (time.Time, error) {
	date := time.Time{}
	err := v.db.Get(&date, getReachingTargetDate, amount, amount)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return date, fmt.Errorf("do get: %w", err)
		}
	}

	return date, nil
}
