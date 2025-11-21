package add

import "github.com/spf13/cobra"

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new object",
	// Run: func(cmd *cobra.Command, args []string) {
	//   fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	// },
}
