package balance

import (
	"context"
	"database/sql"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	"github.com/spf13/cobra"
)

var (
	amount float64
	date   time.Time
)

func init() {
	addBalanceCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "amount")
	addBalanceCmd.Flags().TimeVarP(&date, "date", "d", time.Time{}, []string{time.DateOnly}, "date")
}

var addBalanceCmd = &cobra.Command{
	Use:   "add",
	Short: "add balance record",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {

		app.Queries.CreateBalanceRecord(ctx, db_api.CreateBalanceRecordParams{
			Amount:      amount,
			Date:        date,
			OriginTract: sql.NullInt64{},
		})

		return nil
	}),
}
