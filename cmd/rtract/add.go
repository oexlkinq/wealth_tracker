package rtract

import (
	"context"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/spf13/cobra"
	"github.com/teambition/rrule-go"
)

var (
	desc    string
	amount  float64
	reqsAck bool

	rawFreq  string
	dtstart  time.Time
	interval int
)

func init() {
	addRTractCmd.Flags().StringVarP(&desc, "desc", "d", "", "description")
	addRTractCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "amount")
	addRTractCmd.Flags().BoolVarP(&reqsAck, "reqs-ack", "y", true, "requires acknowledge confirmation?")

	addRTractCmd.Flags().StringVarP(&rawFreq, "freq", "f", "", "frequency")
	addRTractCmd.Flags().TimeVarP(&dtstart, "dtstart", "s", time.Time{}, []string{time.DateOnly}, "start date")
	addRTractCmd.Flags().IntVarP(&interval, "interval", "i", 0, "interval")
}

func init() {
}

var addRTractCmd = &cobra.Command{
	Use:   "add",
	Short: "add balance record",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		freq, err := rrule.StrToFreq(rawFreq)
		if err != nil {
			return err
		}

		r, err := rrule.NewRRule(rrule.ROption{
			Freq:     freq,
			Dtstart:  dtstart,
			Interval: interval,
		})
		if err != nil {
			return err
		}

		return app.Queries.CreateRTract(ctx, db_api.CreateRTractParams{
			Rrule:   r.String(),
			Desc:    desc,
			Amount:  amount,
			ReqsAck: reqsAck,
		})
	}),
}
