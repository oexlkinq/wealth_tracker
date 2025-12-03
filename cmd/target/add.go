package target

import (
	"context"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	"github.com/spf13/cobra"
)

var (
	amount float64
	desc   string
	order  int64
)

func init() {
	addTargetCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "amount")
	addTargetCmd.Flags().StringVarP(&desc, "desc", "d", "", "description")
	addTargetCmd.Flags().Int64VarP(&order, "order", "o", 0, "order")
}

var addTargetCmd = &cobra.Command{
	Use:   "add",
	Short: "add target",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		return app.Queries.CreateTarget(ctx, db_api.CreateTargetParams{
			Amount: amount,
			Desc:   desc,
			Order:  order,
		})
	}),
}
