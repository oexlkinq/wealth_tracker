/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/calc"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wealth_tracker",
	Short: "cli unility to track and calculate reach date of wishlist items",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := app.New()
		if err != nil {
			return fmt.Errorf("create app: %w", err)
		}

		ctx := cmd.Context()

		tx, err := app.DB.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		q := app.Queries.WithTx(tx)

		// удалить сгенерированные транзакции. балансы тоже удалятся изза cascade
		err = q.DeleteGeneratedTracts(ctx)
		if err != nil {
			return fmt.Errorf("delete generated tracts: %w", err)
		}

		// сбор данных для расчёта
		balance, err := q.GetLatestBalanceRecord(ctx)
		if err != nil {
			return fmt.Errorf("get latest balance record: %w", err)
		}

		targets, err := q.ListTargets(ctx)
		if err != nil {
			return fmt.Errorf("list targets: %w", err)
		}

		// расчёт
		tris, err := calc.CalcTargetsReachInfo(ctx, q, balance, targets)
		if err != nil {
			return fmt.Errorf("calc targets reach info: %w", err)
		}

		for _, tri := range tris {
			fmt.Println(tri)
		}

		// завершение
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("commit tx: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
