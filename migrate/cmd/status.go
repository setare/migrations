package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Revert then migrate all migrations",
	Long:  `This command performs a "rewind" command followed by a "migrate".`,
	Run: func(cmd *cobra.Command, args []string) {
		initializePlanner("sql")

		currentMigration, err := target.Current()
		if err != nil {
			output.Error("could get the current: ", err)
			os.Exit(1)
		}

		migrated, err := target.Done()
		if err != nil {
			output.Error("could not list migrations: ", err)
			os.Exit(1)
		}

		output.H1f("Migrations performed: %d total", len(migrated))
		for i, migration := range migrated {
			output.MigrationItemList(i, iconOK(), migration, migration == currentMigration)
		}

		plan, err := planner.Migrate()
		if err != nil {
			output.Error("failed checking pending migrations: ", err)
			os.Exit(1)
		}
		output.Print()
		output.H1f("Migrations pending: %d total", len(plan))
		for i, action := range plan {
			output.PlanAction(i+len(migrated), iconPending(), action)
		}
	},
}

func init() {
	initializeFlags(statusCmd)
	rootCmd.AddCommand(statusCmd)
}
