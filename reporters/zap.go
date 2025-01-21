package reporters

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/jamillosantos/migrations/v2"
)

type runnerReporter struct {
	logger *zap.Logger
}

func NewZapReporter(logger *zap.Logger) migrations.RunnerReporter {
	return &runnerReporter{
		logger: logger,
	}
}

func (r *runnerReporter) BeforeExecute(_ context.Context, req *migrations.BeforeExecuteInfo) {
	if len(req.Plan) == 0 {
		r.logger.Info("no migrations to run")
		return
	}
	ms := make([]string, len(req.Plan))
	for i, action := range req.Plan {
		ms[i] = fmt.Sprintf("%s (%s)", action.Migration.String(), action.Action)
	}
	r.logger.Info(fmt.Sprintf("migration plan with %d migrations", len(req.Plan)), zap.Strings("plan", ms))
}

func (r *runnerReporter) BeforeExecuteMigration(_ context.Context, req *migrations.BeforeExecuteMigrationInfo) {
	r.logger.Info(fmt.Sprintf("migration %s (%s)", req.Migration.String(), req.ActionType))
}

func (r *runnerReporter) AfterExecuteMigration(_ context.Context, req *migrations.AfterExecuteMigrationInfo) {
	if req.Err == nil {
		r.logger.Info(fmt.Sprintf("migration %s (%s) successfully applied", req.Migration.String(), req.ActionType))
		return
	}
	r.logger.Error(fmt.Sprintf("migration %s has failed to execute", req.Migration.String()), zap.Error(req.Err))
}

func (r *runnerReporter) AfterExecute(_ context.Context, req *migrations.AfterExecuteInfo) {
	if req.Err == nil {
		r.logger.Info(fmt.Sprintf("SUCCESS: migration has finished with %d successes and %d failures", len(req.Stats.Successful), len(req.Stats.Errored)))
	} else {
		r.logger.Error(fmt.Sprintf("ERROR: migration has failed with %d successes and %d failures", len(req.Stats.Successful), len(req.Stats.Errored)), zap.Error(req.Err))
	}
	for _, action := range req.Stats.Successful {
		r.logger.Info(fmt.Sprintf("migration %s (%s) was applied", action.Migration.String(), action.Action))
	}
	for _, action := range req.Stats.Errored {
		r.logger.Error(fmt.Sprintf("migration %s (%s) failed to be applied", action.Migration.String(), action.Action))
	}
}
