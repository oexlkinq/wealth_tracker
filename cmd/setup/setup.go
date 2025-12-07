package setup

import (
	"context"
	"fmt"
	"os"

	"github.com/oexlkinq/wealth_tracker/cmd/setup/example"
	"github.com/oexlkinq/wealth_tracker/cmd/setup/models"
	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	SetupCmd.AddCommand(example.ExampleCmd)
}

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "initialize/replace DB by data from .yaml file",
	Args:  cobra.ExactArgs(1),
	RunE: app.MakeCmdRunEFunc(func(cmd *cobra.Command, args []string, ctx context.Context, app *app.App) error {
		data, err := parseDataFile(args[0])
		if err != nil {
			return fmt.Errorf("parse data file: %w", err)
		}

		err = truncateDB(ctx, app.Queries)
		if err != nil {
			return fmt.Errorf("truncate DB: %w", err)
		}
		fmt.Println("DB was truncated")

		err = data.process(ctx, app.Queries)
		if err != nil {
			return fmt.Errorf("process objects: %w", err)
		}

		return nil
	}),
}

type rawData models.RawData

func parseDataFile(filepath string) (*rawData, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	d := yaml.NewDecoder(f)

	var data rawData
	err = d.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func truncateDB(ctx context.Context, qtx *db_api.Queries) error {
	err := qtx.DeleteAllBalanceRecords(ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllRTracts(ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllTargets(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (v *rawData) process(ctx context.Context, qtx *db_api.Queries) error {
	for _, rawBalanceRecord := range v.Balance_records {
		err := rawBalanceRecord.Insert(ctx, qtx)
		if err != nil {
			return err
		}

		fmt.Printf("%+v done\n", rawBalanceRecord)
	}

	for _, rawRtract := range v.Rtracts {
		err := rawRtract.Insert(ctx, qtx)
		if err != nil {
			return err
		}

		fmt.Printf("%+v done\n", rawRtract)
	}

	for _, rawTarget := range v.Targets {
		err := rawTarget.Insert(ctx, qtx)
		if err != nil {
			return err
		}

		fmt.Printf("%+v done\n", rawTarget)
	}

	return nil
}
