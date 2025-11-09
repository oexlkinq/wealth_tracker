package yamldb

import (
	"fmt"
	"log"

	"github.com/oexlkinq/wealth_tracker/internal/models"
	"github.com/teambition/rrule-go"
)

type DB struct {
	Targets        []models.Target
	RTracts        []models.RTract
	BalanceRecords []models.BalanceRecord
}

type rawDB struct {
	Targets        []models.Target
	RTracts        []rawRTract
	BalanceRecords []models.BalanceRecord
}

func (v *rawDB) toDB() *DB {
	db := &DB{
		Targets:        v.Targets,
		RTracts:        make([]models.RTract, len(v.RTracts)),
		BalanceRecords: v.BalanceRecords,
	}

	for i, rawRTract := range v.RTracts {
		res, err := rrule.StrToRRule(rawRTract.RRule)
		if err != nil {
			log.Fatalln(fmt.Errorf("parse rtract rrule: %w", err))
		}

		db.RTracts[i] = models.RTract{
			Amount: rawRTract.Amount,
			Desc:   rawRTract.Desc,
			RRule:  res,
		}
	}

	return db
}

type rawRTract struct {
	Amount float64
	Desc   string
	RRule  string
}

var exampleDB *rawDB = &rawDB{
	Targets: []models.Target{
		models.Target{
			Amount: 12345,
			Order:  0,
			Desc:   "какаято цель",
		},
		models.Target{
			Amount: 2345,
			Order:  1,
			Desc:   "ещё какаято цель",
		},
	},
	RTracts: []rawRTract{
		rawRTract{
			Amount: 6000,
			Desc:   "зп каждый месяц",
			RRule:  "",
		},
	},
}
