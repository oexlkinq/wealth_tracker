package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/spf13/cobra"
)

var listBalancesCmd = &cobra.Command{
	Use:   "list",
	Short: "list balance records",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		brs, err := app.Queries.ListBalanceRecords(ctx)
		if err != nil {
			return err
		}

		for _, br := range brs {
			fmt.Printf("{id: %5d, date: %s, amount: %10.2f}\n", br.ID, br.Date.Format(time.DateOnly), br.Amount)
		}

		return nil
	}),
}
