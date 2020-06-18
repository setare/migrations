package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ory/viper"
	"github.com/prometheus/common/log"
	"github.com/setare/migrations"
	"github.com/spf13/cobra"
)

var (
	Target           migrations.Target
	Source           migrations.Source
	planner          *migrations.Planner
	Runner           *migrations.Runner
	ExecutionContext migrations.ExecutionContext
)

func Initialize(source migrations.Source, target migrations.Target, executionContext migrations.ExecutionContext) error {
	Runner = migrations.NewRunner(source, target)
	Target, Source = target, source
	ExecutionContext = executionContext

	err := Target.Create()
	if err != nil {
		return err
	}

	return nil
}

func initializeFlags(cmd *cobra.Command) {
	var planFlag bool
	cmd.Flags().BoolVarP(&planFlag, "plan", "p", false, "When active, informs commands to never run, only plan their actions.")
	viper.BindPFlag("plan", cmd.Flags().Lookup("plan"))
}

func PlanAndRun(plannerFunc migrations.PlannerFunc) {
	plan, err := plannerFunc(Source, Target).Plan()
	if err != nil {
		log.Error("planning failed: ", err)
		os.Exit(800)
	}

	reportPlan(plan)

	if viper.GetBool("plan") || len(plan) == 0 {
		return
	}

	if !viper.GetBool("yes") {
		proceed := false
		prompt := &survey.Confirm{
			Message: "Proceed with the plan?",
		}
		survey.AskOne(prompt, &proceed)
		if !proceed {
			Output.Warn("user cancelled")
			os.Exit(EC_CANCELLED)
		}
	}

	Output.H1("Execution")
	stats, err := Runner.Execute(ExecutionContext, plan, &Output)
	Output.Print()
	Output.Printf("%s successful, %s errors", StyleSuccess.Sprint(len(stats.Successful)), StyleError.Sprint(len(stats.Errored)))
	Output.Print()
	if err != nil {
		Output.Error(err)
		os.Exit(500)
	}
}
