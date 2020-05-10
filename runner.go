package migrations

import "github.com/pkg/errors"

// Runner will receive the `Plan` from the `Planner` and execute it.
type Runner struct {
	source Source
	target Target
}

type RunnerOption func(*Runner)

func NewRunner(source Source, target Target) *Runner {
	return &Runner{
		source: source,
		target: target,
	}
}

type RunnerReporter interface {
	BeforeExecute(ActionType, Migration)
	AfterExecute(ActionType, Migration, error)
}

type ExecutionStats struct {
	Successful []*Action
	Errored    []*Action
}

func (runner *Runner) Execute(executionContext ExecutionContext, plan Plan, reporter RunnerReporter) (*ExecutionStats, error) {
	stats := &ExecutionStats{
		Successful: make([]*Action, 0, len(plan)),
	}
	for _, action := range plan {
		var err error
		reporter.BeforeExecute(action.Action, action.Migration)
		switch action.Action {
		case ActionTypeDo:
			err = action.Migration.Do(executionContext)
			if err == nil {
				runner.target.Add(action.Migration)
			}
		case ActionTypeUndo:
			if !action.Migration.CanUndo() {
				err = WrapMigration(ErrMigrationNotUndoable, action.Migration)
			} else {
				err = action.Migration.Undo(executionContext)
				if err == nil {
					runner.target.Remove(action.Migration)
				}
			}
		default:
			err = errors.Wrap(ErrInvalidAction, string(action.Action))
		}
		reporter.AfterExecute(action.Action, action.Migration, err)
		if err != nil {
			stats.Errored = []*Action{action}
			return stats, err
		}
	}
	return stats, nil
}
