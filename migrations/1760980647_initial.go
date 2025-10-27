package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		return app.RunInTransaction(func(txApp core.App) error {
			budget_flows := core.NewBaseCollection("budget_flows")
			budget_flows.ListRule = types.Pointer("")
			budget_flows.ViewRule = types.Pointer("")
			budget_flows.CreateRule = types.Pointer("")
			budget_flows.UpdateRule = types.Pointer("")
			budget_flows.DeleteRule = types.Pointer("")
			budget_flows.Fields.Add(
				&core.TextField{
					Name:     "desc",
					Required: true,
				},
				&core.TextField{
					Name:     "rrule",
					Required: true,
				},
				&core.NumberField{
					Name:     "amount",
					Required: true,
				},
			)
			err := app.Save(budget_flows)
			if err != nil {
				return err
			}

			targets := core.NewBaseCollection("targets")
			targets.ListRule = types.Pointer("")
			targets.ViewRule = types.Pointer("")
			targets.CreateRule = types.Pointer("")
			targets.UpdateRule = types.Pointer("")
			targets.DeleteRule = types.Pointer("")
			targets.Fields.Add(
				&core.TextField{
					Name:     "desc",
					Required: true,
				},
				&core.NumberField{
					Name:     "amount",
					Required: true,
				},
				&core.NumberField{
					Name:     "order",
					Required: true,
				},
			)
			err = app.Save(targets)
			if err != nil {
				return err
			}

			budget := core.NewBaseCollection("budget")
			budget.ListRule = types.Pointer("")
			budget.ViewRule = types.Pointer("")
			budget.CreateRule = types.Pointer("")
			budget.UpdateRule = types.Pointer("")
			budget.DeleteRule = types.Pointer("")
			budget.Fields.Add(
				&core.NumberField{
					Name:     "amount",
					Required: true,
				},
				&core.AutodateField{
					Name:     "date",
					OnCreate: true,
				},
			)
			err = app.Save(budget)
			if err != nil {
				return err
			}

			return nil
		})
	}, func(app core.App) error {
		return app.RunInTransaction(func(txApp core.App) error {
			for _, collectionName := range []string{"budget_flows", "targets", "budget"} {
				collection, err := app.FindCollectionByNameOrId(collectionName)
				if err != nil {
					return err
				}

				err = app.Delete(collection)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
}
