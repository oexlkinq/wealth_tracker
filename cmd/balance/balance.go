package balance

import (
	"github.com/spf13/cobra"
)

func init() {
	BalanceCmd.AddCommand(addBalanceCmd)
	BalanceCmd.AddCommand(listBalancesCmd)
}

var BalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "CRUD balance records",
}
