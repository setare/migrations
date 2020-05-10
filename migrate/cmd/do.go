package cmd

import (
	"github.com/setare/migrations"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Steps one migration forward",
	Long:  `Starting from the current migration, this command will do one.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")
		planAndRun(func() (migrations.Plan, error) {
			return planner.Step(1)
		})
	},
}

func init() {
	initializeFlags(doCmd)
	rootCmd.AddCommand(doCmd)
}
