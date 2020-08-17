package status

import (
	"errors"
	"os"

	"github.com/jamillosantos/migrations"
	"github.com/jamillosantos/migrations/cmd"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Revert then migrate all migrations",
	Long:  `This command performs a "rewind" command followed by a "migrate".`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		currentMigration, err := cmd.Target.Current()
		if err != nil && !errors.Is(err, migrations.ErrNoCurrentMigration) {
			cmd.Output.Error("failed getting the current migration: ", err)
			os.Exit(1)
		}

		migrated, err := cmd.Target.Done()
		if err != nil {
			cmd.Output.Error("could not list migrations: ", err)
			os.Exit(1)
		}

		cmd.Output.H1f("Migrations performed: %d total", len(migrated))
		for i, migration := range migrated {
			cmd.Output.MigrationItemList(i, cmd.IconOK(), migration, migration == currentMigration)
		}

		plan, err := migrations.MigratePlanner(cmd.Source, cmd.Target).Plan()
		if err != nil {
			cmd.Output.Error("failed checking pending migrations: ", err)
			os.Exit(1)
		}
		cmd.Output.Print()
		cmd.Output.H1f("Migrations pending: %d total", len(plan))
		for i, action := range plan {
			cmd.Output.PlanAction(i+len(migrated), cmd.IconPending(), action)
		}
	},
}
