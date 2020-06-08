package reset

import (
	"github.com/setare/migrations/cmd"
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Revert then migrate all migrations",
	Long:  `This command performs a "rewind" command followed by a "migrate".`,
	Run: func(_ *cobra.Command, args []string) {
		cmd.PlanAndRun(cmd.Planner.Reset)
	},
}
