package example

import (
	"math/rand/v2"
	"os"
	"time"

	"github.com/oexlkinq/wealth_tracker/cmd/setup/models"
	"github.com/spf13/cobra"
	"github.com/teambition/rrule-go"
	"gopkg.in/yaml.v3"
)

var ExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "print example data file",
	Run: func(cmd *cobra.Command, args []string) {
		enc := yaml.NewEncoder(os.Stdout)

		data := &models.RawData{
			Balance_records: []*models.RawBalanceRecord{
				{
					Amount: 50_000,
					Date:   "2025-12-06",
				},
			},
			Rtracts: []*models.RawRTract{
				{
					Rrule:    randRRule(),
					Desc:     "зп",
					Amount:   27_500,
					Reqs_ack: false,
				},
				{
					Rrule:    randRRule(),
					Desc:     "аренда квартиры",
					Amount:   -15_000,
					Reqs_ack: false,
				},
				{
					Rrule:    randRRule(),
					Desc:     "за инет",
					Amount:   -890,
					Reqs_ack: true,
				},
			},
			Targets: []*models.RawTarget{
				{
					Amount: 50_000,
					Desc:   "на чёрный день",
					Order:  0,
				},
				{
					Amount: 55_000,
					Desc:   "на время без работы",
					Order:  0,
				},
				{
					Amount: 150,
					Desc:   "носки на нг",
					Order:  1,
				},
				{
					Amount: 50_000,
					Desc:   "на компик",
					Order:  2,
				},
			},
		}

		err := enc.Encode(data)
		if err != nil {
			panic(err)
		}
	},
}

const maxDaysShift = 30

func randRRule() string {
	// смещение в днях. случайное из [-maxDaysShift, 0]
	daysShift := -maxDaysShift + rand.IntN(maxDaysShift+1)

	r, err := rrule.NewRRule(rrule.ROption{
		Freq:     rrule.MONTHLY,
		Dtstart:  time.Now().Add(time.Hour * 24 * time.Duration(daysShift)),
		Interval: 1,
	})
	if err != nil {
		panic(err)
	}

	return r.String()
}
