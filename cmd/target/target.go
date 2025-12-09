package target

import (
	"github.com/spf13/cobra"
)

func init() {
	TargetCmd.AddCommand(addTargetCmd)
	TargetCmd.AddCommand(listTargetsCmd)
}

var TargetCmd = &cobra.Command{
	Use:   "target",
	Short: "CRUD targets",
}
