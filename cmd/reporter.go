package cmd

import "github.com/setare/migrations"

func reportPlan(plan migrations.Plan) {
	if len(plan) == 0 {
		Output.Success("no operation needed")
		return
	}
	Output.H1("The Plan")
	for i, action := range plan {
		Output.PlanAction(i, "", action)
	}
	Output.Print()
}
