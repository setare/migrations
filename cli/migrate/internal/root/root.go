package root

import (
	"fmt"
	"os"

	"github.com/jamillosantos/migrations/cli/migrate/internal/cmdsql"
	"github.com/jamillosantos/migrations/cmd"
	"github.com/jamillosantos/migrations/cmd/create"
	"github.com/jamillosantos/migrations/cmd/do"
	"github.com/jamillosantos/migrations/cmd/migrate"
	"github.com/jamillosantos/migrations/cmd/reset"
	"github.com/jamillosantos/migrations/cmd/rewind"
	"github.com/jamillosantos/migrations/cmd/status"
	"github.com/jamillosantos/migrations/cmd/undo"
	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   fmt.Sprint(os.Args[0]),
	Short: "The migration CLI tool",
	Long:  `migrate is a CLI tool for Go that let the developer migrate anything.`,
	PersistentPreRun: func(_ *cobra.Command, args []string) {
		source, target, err := cmdsql.Initialize()
		if err != nil {
			cmd.Output.Error(err)
			os.Exit(1)
		}
		err = cmd.Initialize(source, target, cmdsql.NewExecutionContext())
		if err != nil {
			cmd.Output.Error(err)
			os.Exit(2)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.Output.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	var (
		configFileFlag string
		yesFlag        bool
	)

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", "migrations.yaml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Disable confirmations to all confirmations before executing any migration.")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))

	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(migrate.MigrateCmd)
	rootCmd.AddCommand(rewind.RewindCmd)
	rootCmd.AddCommand(reset.ResetCmd)
	rootCmd.AddCommand(do.DoCmd)
	rootCmd.AddCommand(undo.UndoCmd)
	rootCmd.AddCommand(status.StatusCmd)
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Use config file from the flag.
	if cfgFile := viper.GetString("config"); cfgFile != "" {
		cmd.Output.Warn("using config: ", cfgFile)
		viper.SetConfigFile(viper.GetString("config"))
	}
	viper.BindEnv("directory")
	viper.BindEnv("dsn")
	viper.BindEnv("driver")

	if err := viper.ReadInConfig(); err != nil {
		cmd.Output.Errorf("failed loading %s:", viper.ConfigFileUsed())
		os.Exit(999)
	}
}
