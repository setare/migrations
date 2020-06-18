package do

import (
	"github.com/setare/migrations"
	"github.com/setare/migrations/cmd"
	"github.com/spf13/cobra"
)

var DoCmd = &cobra.Command{
	Use:   "do",
	Short: "Steps one migration forward",
	Long:  `Starting from the current migration, this command will do one.`,
	Run: func(_ *cobra.Command, args []string) {
		cmd.PlanAndRun(migrations.DoPlanner)
	},
}
