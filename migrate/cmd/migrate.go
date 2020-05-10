package cmd

import (
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate <type>",
	Short: "Run all pending migrations",
	Long: `This command will check the current migration and run all migrations
that are listed after it. The migrations will run in order and if one fail, the
process will be canceled. After each migration successfully ran, the system will
save the state.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")

		planAndRun(planner.Migrate)
	},
}

func init() {
	initializeFlags(migrateCmd)
	rootCmd.AddCommand(migrateCmd)
}
