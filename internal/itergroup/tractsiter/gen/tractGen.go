package gen

import (
	"fmt"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/models"
	"github.com/teambition/rrule-go"
)

// стейт генератора дат ртракта
type GTractGen struct {
	*db_api.Rtract
	rr   *rrule.RRule
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
		rr:     rr,
		Next:   rr.Iterator(),
	}, nil
}

// сменить дату начала генерации и пересоздать итератор Next. inc определяет будет ли использована дата dt, если она так же является вхождением rrule
func (v *GTractGen) ResetTo(dt time.Time, inc bool) {
	v.rr.DTStart(v.rr.After(dt, inc))
	v.Next = v.rr.Iterator()
}

func (v *GTractGen) All() iter.Seq[*models.CalcTract] {
	return func(yield func(*models.CalcTract) bool) {
		for {
			date, ok := v.Next()
			if !ok {
				return
			}
			fmt.Printf("gen %s %s\n", v.Rtract.Desc, date)

			if !yield(&models.CalcTract{
				Tract: &db_api.Tract{
					ID:     -1,
					Amount: v.Rtract.Amount,
					Date:   date,
					Type:   "rtract",
					Acked:  false,
				},
				RTractID:  v.ID,
				Generated: true,
			}) {
				return
			}
		}
	}
}
