package root

import (
	"os"

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

var RootCmd = &cobra.Command{
	Use:   "migrations",
	Short: "The migration CLI tool",
	Long:  `migrate is a CLI tool for Go that let the developer migrate anything.`,
}

func Init() {
	cobra.OnInitialize(initConfig)

	var (
		configFileFlag string
		yesFlag        bool
	)

	RootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", "migrations.yaml", "config file")
	RootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Disable confirmations to all confirmations before executing any migration.")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("yes", RootCmd.PersistentFlags().Lookup("yes"))

	RootCmd.AddCommand(create.CreateCmd)
	RootCmd.AddCommand(migrate.MigrateCmd)
	RootCmd.AddCommand(rewind.RewindCmd)
	RootCmd.AddCommand(reset.ResetCmd)
	RootCmd.AddCommand(do.DoCmd)
	RootCmd.AddCommand(undo.UndoCmd)
	RootCmd.AddCommand(status.StatusCmd)
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
