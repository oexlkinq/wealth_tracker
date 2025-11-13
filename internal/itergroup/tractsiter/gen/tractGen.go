package gen

import (
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	"github.com/teambition/rrule-go"
)

// стейт генератора дат ртракта
type GTractGen struct {
	*db_api.Rtract
	Next rrule.Next
}

func New(rtract *db_api.Rtract, since time.Time) (*GTractGen, error) {
	rr, err := rrule.StrToRRule(rtract.Rrule)
	if err != nil {
		return nil, err
	}
	rr.DTStart(rr.After(since, true))

	return &GTractGen{
		Rtract: rtract,
		Next:   rr.Iterator(),
	}, nil
}

func (v *GTractGen) All() iter.Seq[*db_api.ListTractsSinceRow] {
	return func(yield func(*db_api.ListTractsSinceRow) bool) {
		for {
			date, ok := v.Next()
			if !ok {
				return
			}

			if !yield(&db_api.ListTractsSinceRow{
				Amount: v.Rtract.Amount,
				Date:   date,
			}) {
				return
			}
		}
	}
}
