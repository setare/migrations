package zap

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	migrations "github.com/jamillosantos/migrations"
)

type runnerReporter struct {
	logger *zap.Logger
}

func NewRunnerReport(logger *zap.Logger) migrations.RunnerReporter {
	return &runnerReporter{
		logger: logger,
	}
}

func (r *runnerReporter) BeforeExecute(_ context.Context, plan migrations.Plan) {
	if len(plan) == 0 {
		r.logger.Info("no migrations to run")
		return
	}
	ms := make([]string, len(plan))
	for i, action := range plan {
		ms[i] = fmt.Sprintf("%s (%s)", action.Migration.String(), action.Action)
	}
	r.logger.Info(fmt.Sprintf("migration plan with %d migrations", len(plan)), zap.Strings("plan", ms))
}

func (r *runnerReporter) BeforeExecuteMigration(ctx context.Context, actionType migrations.ActionType, migration migrations.Migration) {
	r.logger.Info(fmt.Sprintf("migration %s (%s)", migration.String(), actionType))
}

func (r *runnerReporter) AfterExecuteMigration(ctx context.Context, actionType migrations.ActionType, migration migrations.Migration, err error) {
	if err == nil {
		r.logger.Info(fmt.Sprintf("migration %s (%s) successfully applied", migration.String(), actionType))
		return
	}
	r.logger.Error(fmt.Sprintf("migration %s has failed to execute", migration.String()), zap.Error(err))
}

func (r *runnerReporter) AfterExecute(_ context.Context, _ migrations.Plan, stats *migrations.ExecutionStats, err error) {
	if err == nil {
		r.logger.Info(fmt.Sprintf("SUCCESS: migration has finished with %d successes and %d failures", len(stats.Successful), len(stats.Errored)))
	} else {
		r.logger.Error(fmt.Sprintf("ERROR: migration has failed with %d successes and %d failures", len(stats.Successful), len(stats.Errored)), zap.Error(err))
	}
	for _, action := range stats.Successful {
		r.logger.Info(fmt.Sprintf("migration %s (%s) was applied", action.Migration.String(), action.Action))
	}
	for _, action := range stats.Errored {
		r.logger.Error(fmt.Sprintf("migration %s (%s) failed to be applied", action.Migration.String(), action.Action))
	}
}
