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

// Execute performs a plan, running all actions migration by migration.
//
// Before running, Execute will check for Undo actions that cannot be performmed into undoable migrations. If that
// happens, an `ErrMigrationNotUndoable` will be returned and nothing will be executed.
//
// For each migration executed, the system will move the cursor to that point. So that, if any error happens during the
// migration execution (do or undo), the execution will be stopped and the error will be returned. All performed actions
// WILL NOT be rolled back.
func (runner *Runner) Execute(executionContext ExecutionContext, plan Plan, reporter RunnerReporter) (*ExecutionStats, error) {
	stats := &ExecutionStats{
		Successful: make([]*Action, 0, len(plan)),
	}

	// Check for undoable migrations...
	for _, action := range plan {
		if action.Action == ActionTypeUndo && !action.Migration.CanUndo() {
			return nil, WrapMigration(ErrMigrationNotUndoable, action.Migration)
		}
	}

	for _, action := range plan {
		var err error
		if reporter != nil {
			reporter.BeforeExecute(action.Action, action.Migration)
		}
		switch action.Action {
		case ActionTypeDo:
			err = action.Migration.Do(executionContext)
			if err == nil {
				runner.target.Add(action.Migration)
			}
		case ActionTypeUndo:
			// Undoable migrations were already checked before.
			err = action.Migration.Undo(executionContext)
			if err == nil {
				runner.target.Remove(action.Migration)
			}
		default:
			err = errors.Wrap(ErrInvalidAction, string(action.Action))
		}
		if reporter != nil {
			reporter.AfterExecute(action.Action, action.Migration, err)
		}
		if err == nil {
			stats.Successful = append(stats.Successful, action)
		} else {
			stats.Errored = []*Action{action}
			return stats, err
		}
	}
	return stats, nil
}
