package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/prometheus/common/log"
	"github.com/setare/migrations"
	"github.com/setare/migrations/migrate/cmd/cmdsql"
	"github.com/setare/migrations/migrate/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	source migrations.Source
	target migrations.Target
)

func initializePlanner(sourceType string) {
	// TODO(Jota): Move this to the specific command

	switch sourceType {
	case "sql":
		if !viper.IsSet("dsn") || viper.Get("dsn") == "" {
			output.Error("--dsn or DSN environment variable not defined")
			os.Exit(2)
		}

		err := utils.Spin(func(s *spinner.Spinner) error {
			s.Suffix = "Connecting ..."
			return cmdsql.Connect(viper.GetString("driver"), viper.GetString("dsn"))
		})
		if err != nil {
			output.Error(err)
			os.Exit(3)
		}

		source, target, err = cmdsql.Initialize(viper.GetString("directory"))
		if err != nil {
			output.Error("failed initializing source or target: ", err)
			os.Exit(4)
		}
	default:
		output.Error("unknown type: ", sourceType)
		os.Exit(5)
	}
	planner = migrations.NewPlanner(source, target)

	runner = migrations.NewRunner(source, target)
}

func initializeFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&planFlag, "plan", "p", false, "When active, informs commands to never run, only plan their actions.")
}

func planAndRun(action func() (migrations.Plan, error)) {
	plan, err := action()
	if err != nil {
		log.Error("planning failed: ", err)
		os.Exit(800)
	}

	reportPlan(plan)

	if planFlag || len(plan) == 0 {
		return
	}

	if !yesFlag {
		proceed := false
		prompt := &survey.Confirm{
			Message: "Proceed with the plan?",
		}
		survey.AskOne(prompt, &proceed)
		if !proceed {
			output.Warn("user cancelled")
			os.Exit(EC_CANCELLED)
		}
	}

	output.H1("Execution")
	stats, err := runner.Execute(cmdsql.NewExecutionContext(), plan, &output)
	output.Print()
	output.Printf("%s successful, %s errors", styleSuccess.Sprint(len(stats.Successful)), styleError.Sprint(len(stats.Errored)))
	output.Print()
	if err != nil {
		output.Error(err)
		os.Exit(500)
	}
}
