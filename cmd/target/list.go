package target

import (
	"context"
	"fmt"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/spf13/cobra"
)

var listTargetsCmd = &cobra.Command{
	Use:   "list",
	Short: "list targets",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		targets, err := app.Queries.ListTargets(ctx)
		if err != nil {
			return err
		}

		for _, target := range targets {
			fmt.Printf("{id: %5d, amount: %10.2f, order: %3d, desc: %s}\n", target.ID, target.Amount, target.Order, target.Desc)
		}

		return nil
	}),
}
