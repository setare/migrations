package migrations

import (
	"context"
)

func Migrate(ctx context.Context, runner *Runner, runnerReporter RunnerReporter) (*ExecutionStats, error) {
	err := runner.target.Create()
	if err != nil {
		return nil, err
	}

	plan, err := MigratePlanner(runner.source, runner.target).Plan()
	if err != nil {
		return nil, err
	}
	return runner.Execute(ctx, plan, runnerReporter)
}
