package cmd

import (
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Revert then migrate all migrations",
	Long:  `This command performs a "rewind" command followed by a "migrate".`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")
		planAndRun(planner.Reset)
	},
}

func init() {
	initializeFlags(resetCmd)
	rootCmd.AddCommand(resetCmd)
}
