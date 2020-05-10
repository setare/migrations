package cmd

import (
	"github.com/spf13/cobra"
)

var rewindCmd = &cobra.Command{
	Use:   "rewind",
	Short: "Undo all migrations",
	Long: `Starting from the current migration, this command will undo all
migrations from the most recent to the first one.

Some migrations cannot be undone, if one of those are found the process will
undo all migrations until fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")
		planAndRun(planner.Rewind)
	},
}

func init() {
	initializeFlags(rewindCmd)
	rootCmd.AddCommand(rewindCmd)
}
