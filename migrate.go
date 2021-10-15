package migrations

import (
	"context"
)

func Migrate(ctx context.Context, runner *Runner, runnerReporter RunnerReporter) (*ExecutionStats, error) {
	plan, err := MigratePlanner(runner.source, runner.target).Plan()
	if err != nil {
		return nil, err
	}
	return runner.Execute(ctx, plan, runnerReporter)
}
