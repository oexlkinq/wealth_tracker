package rtract

import (
	"github.com/spf13/cobra"
)

func init() {
	RtractCmd.AddCommand(addRTractCmd)
	RtractCmd.AddCommand(listRTractsCmd)
}

var RtractCmd = &cobra.Command{
	Use:   "rtract",
	Short: "CRUD rtracts",
}
