package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   fmt.Sprint(os.Args[0]),
	Short: "The migration CLI tool",
	Long:  `migrate is a CLI tool for Go that let the developer migrate anything.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", "migrations.yaml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Disable confirmations to all confirmations before executing any migration.")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Use config file from the flag.
	viper.SetConfigFile(configFileFlag)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		output.Errorf("failed loading %s:", viper.ConfigFileUsed())
		os.Exit(999)
	}
}
