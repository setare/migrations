package undo

import (
	"github.com/setare/migrations"
	"github.com/setare/migrations/cmd"
	"github.com/spf13/cobra"
)

var UndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the most recent migration",
	Long: `Some migrations cannot be undone, if one of those are found the process will
undo all migrations until fail.`,
	Run: func(_ *cobra.Command, args []string) {
		cmd.PlanAndRun(func() (migrations.Plan, error) {
			return cmd.Planner.Step(1)
		})
	},
}
