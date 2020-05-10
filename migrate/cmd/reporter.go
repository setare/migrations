package cmd

import "github.com/setare/migrations"

func reportPlan(plan migrations.Plan) {
	if len(plan) == 0 {
		output.Success("no operation needed")
		return
	}
	output.H1("The Plan")
	for i, action := range plan {
		output.PlanAction(i, "", action)
	}
	output.Print()
}
