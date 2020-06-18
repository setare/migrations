package rewind

import (
	"github.com/setare/migrations"
	"github.com/setare/migrations/cmd"
	"github.com/spf13/cobra"
)

var RewindCmd = &cobra.Command{
	Use:   "rewind",
	Short: "Undo all migrations",
	Long: `Starting from the current migration, this command will undo all
migrations from the most recent to the first one.

Some migrations cannot be undone, if one of those are found the process will
undo all migrations until fail.`,
	Run: func(_ *cobra.Command, args []string) {
		cmd.PlanAndRun(migrations.RewindPlanner)
	},
}
