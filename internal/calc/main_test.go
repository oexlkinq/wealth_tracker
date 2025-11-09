package calc_test

import (
	"testing"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/calc"
	"github.com/oexlkinq/wealth_tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
)

func TestCalcTargetsReachInfo(t *testing.T) {
	today := time.Now().Round(time.Hour * 24)
	rtracts := []*models.RTract{
		&models.RTract{
			Amount: 15000,
			Desc:   "зп",
			RRule: makeRRule(rrule.ROption{
				Dtstart:  today,
				Freq:     rrule.MONTHLY,
				Interval: 1,
			}),
		},
		&models.RTract{
			Amount: -10000,
			Desc:   "трата",
			RRule: makeRRule(rrule.ROption{
				Dtstart:  today.Add(time.Hour * 24 * 14),
				Freq:     rrule.MONTHLY,
				Interval: 1,
			}),
		},
	}

	latestBalanceRecord := models.BalanceRecord{
		Amount: 10000,
		Date:   time.Now().Add(-time.Hour * 24),
	}

	targets := []models.Target{
		models.Target{
			Amount: 35000,
			Order:  0,
			Desc:   "на чёрный день",
		},
		models.Target{
			Amount: 15000,
			Order:  1,
			Desc:   "рандомная цель",
		},
		models.Target{
			Amount: 1_000_000_000,
			Order:  1,
			Desc:   "цель гигачад",
		},
	}

	_, ok := calc.CalcTargetsReachInfo(rtracts, today, latestBalanceRecord, targets)
	assert.False(t, ok)

	// TODO: доделать тест
}

func makeRRule(opts rrule.ROption) *rrule.RRule {
	res, err := rrule.NewRRule(opts)
	if err != nil {
		panic(err)
	}

	return res
}
