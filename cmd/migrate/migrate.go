package migrate

import (
	"github.com/setare/migrations"
	"github.com/setare/migrations/cmd"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate <type>",
	Short: "Run all pending migrations",
	Long: `This command will check the current migration and run all migrations
that are listed after it. The migrations will run in order and if one fail, the
process will be canceled. After each migration successfully ran, the system will
save the state.`,
	Run: func(_ *cobra.Command, args []string) {
		cmd.PlanAndRun(migrations.MigratePlanner)
	},
}
