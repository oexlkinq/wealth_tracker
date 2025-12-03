package cmd

import (
	"context"
	"fmt"

	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "just for tests",
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		fmt.Println("hii")

		return nil
	}),
}
