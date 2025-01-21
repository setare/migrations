package migrations

import (
	"context"
	"fmt"
)

// Runner will receive the `Plan` from the `Planner` and execute it.
type Runner struct {
	reporter RunnerReporter
	source   Source
	target   Target
}

type runnerOptions struct {
	Reporter RunnerReporter
}

type RunnerOption func(*runnerOptions)

func NewRunner(source Source, target Target, options ...RunnerOption) *Runner {
	opts := runnerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return &Runner{
		source:   source,
		target:   target,
		reporter: opts.Reporter,
	}
}

func WithReporter(reporter RunnerReporter) RunnerOption {
	return func(options *runnerOptions) {
		options.Reporter = reporter
	}
}

type BeforeExecuteInfo struct {
	Plan Plan
}

type BeforeExecuteMigrationInfo struct {
	ActionType ActionType
	Migration  Migration
}

type AfterExecuteMigrationInfo struct {
	ActionType ActionType
	Migration  Migration
	Err        error
}

type AfterExecuteInfo struct {
	Plan  Plan
	Stats *ExecutionResponse
	Err   error
}

type RunnerReporter interface {
	BeforeExecute(ctx context.Context, info *BeforeExecuteInfo)
	BeforeExecuteMigration(ctx context.Context, info *BeforeExecuteMigrationInfo)
	AfterExecuteMigration(ctx context.Context, info *AfterExecuteMigrationInfo)
	AfterExecute(ctx context.Context, info *AfterExecuteInfo)
}

type ExecutionResponse struct {
	Successful []*Action
	Errored    []*Action
}

type ExecuteRequest struct {
	Plan Plan
}

// Execute performs a plan, running all actions migration by migration.
//
// Before running, Execute will check for Undo actions that cannot be performed into undoable migrations. If that
// happens, an `ErrMigrationNotUndoable` will be returned and nothing will be executed.
//
// For each migration executed, the system will move the cursor to that point. So that, if any error happens during the
// migration execution (do or undo), the execution will be stopped and the error will be returned. All performed actions
// WILL NOT be rolled back.
func (runner *Runner) Execute(ctx context.Context, req *ExecuteRequest) (ExecutionResponse, error) {
	stats := ExecutionResponse{
		Successful: make([]*Action, 0, len(req.Plan)),
	}

	// Check for undoable migrations...
	for _, action := range req.Plan {
		if action.Action == ActionTypeUndo && !action.Migration.CanUndo() {
			return stats, WrapMigration(ErrMigrationNotUndoable, action.Migration)
		}
	}

	if runner.reporter != nil {
		runner.reporter.BeforeExecute(ctx, &BeforeExecuteInfo{
			Plan: req.Plan,
		})
	}

	for _, action := range req.Plan {
		var err error
		if runner.reporter != nil {
			runner.reporter.BeforeExecuteMigration(ctx, &BeforeExecuteMigrationInfo{
				ActionType: action.Action,
				Migration:  action.Migration,
			})
		}
		switch action.Action {
		case ActionTypeDo:
			err = runner.target.Add(ctx, action.Migration.ID())
			if err != nil {
				return stats, err
			}
			err = action.Migration.Do(ctx)
			if err == nil {
				err = runner.target.FinishMigration(ctx, action.Migration.ID())
			}
		case ActionTypeUndo:
			// Undoable migrations were already checked before.
			err = runner.target.StartMigration(ctx, action.Migration.ID())
			if err != nil {
				return ExecutionResponse{}, err
			}
			err = action.Migration.Undo(ctx)
			if err == nil {
				err = runner.target.Remove(ctx, action.Migration.ID())
			}
		default:
			err = fmt.Errorf("%w: %s", ErrInvalidAction, string(action.Action))
		}
		if runner.reporter != nil {
			runner.reporter.AfterExecuteMigration(ctx, &AfterExecuteMigrationInfo{
				ActionType: action.Action,
				Migration:  action.Migration,
				Err:        err,
			})
		}
		if err == nil {
			stats.Successful = append(stats.Successful, action)
		} else {
			stats.Errored = []*Action{action}

			if runner.reporter != nil {
				runner.reporter.AfterExecute(ctx, &AfterExecuteInfo{
					Plan:  req.Plan,
					Stats: &stats,
					Err:   err,
				})
			}

			return stats, err
		}
	}
	if runner.reporter != nil {
		runner.reporter.AfterExecute(ctx, &AfterExecuteInfo{
			Plan:  req.Plan,
			Stats: &stats,
			Err:   nil,
		})
	}
	return stats, nil
}
