package migrations

import (
	"context"
)

func Migrate(ctx context.Context, runner *Runner, runnerReporter RunnerReporter) (*ExecutionStats, error) {
	if locker, ok := runner.target.(TargetLocker); ok {
		lock, err := locker.Lock()
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = lock.Unlock()
		}()
	}

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
