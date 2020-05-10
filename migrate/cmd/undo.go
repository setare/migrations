package cmd

import (
	"github.com/setare/migrations"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the most recent migration",
	Long: `Some migrations cannot be undone, if one of those are found the process will
undo all migrations until fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")
		planAndRun(func() (migrations.Plan, error) {
			return planner.Step(1)
		})
	},
}

func init() {
	initializeFlags(undoCmd)
	rootCmd.AddCommand(undoCmd)
}
