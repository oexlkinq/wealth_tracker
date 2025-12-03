package rtract

import (
	"context"
	"fmt"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/spf13/cobra"
	"github.com/teambition/rrule-go"
)

var listRTractsCmd = &cobra.Command{
	Use:   "list",
	Short: "list rtracts",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		targets, err := app.Queries.ListRTracts(ctx)
		if err != nil {
			return err
		}

		for _, target := range targets {
			r, err := rrule.StrToRRule(target.Rrule)
			if err != nil {
				panic(fmt.Errorf("bad rrule from db (id of target: %d)", target.ID))
			}

			nextTime := r.After(time.Now(), true)

			fmt.Printf("{id: %5d, amount: %10.2f, next: %s, desc: %s}\n", target.ID, target.Amount, nextTime, target.Desc)
		}

		return nil
	}),
}
