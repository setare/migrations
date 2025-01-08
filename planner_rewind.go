package migrations

import (
	"context"
	"errors"
)

type rewindPlanner struct {
	source Source
	target Target
}

// RewindPlanner is a planner that will undo all migrations starting from the current revision.
func RewindPlanner(source Source, target Target) Planner {
	return &rewindPlanner{
		source: source,
		target: target,
	}
}

func (planner *rewindPlanner) Plan(ctx context.Context) (Plan, error) {
	repo, err := planner.source.Load(ctx)
	if err != nil {
		return nil, err
	}

	currentMigrationID, err := planner.target.Current(ctx)
	if errors.Is(err, ErrNoCurrentMigration) {
		// If there is no current migration, no migrations should run.
		return Plan{}, nil
	} else if err != nil {
		// Otherwise, it is just an error.
		return nil, err
	}

	migrationList, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}

	currentMigrationIndex, err := findMigrationIndexByID(migrationList, currentMigrationID)
	if err != nil {
		return nil, err
	}

	// Build the plan
	lst := migrationList[:currentMigrationIndex+1]
	plan := make(Plan, len(lst))
	for i, m := range lst {
		// Inverts the order of the list, the rewind should be planned in the inverse execution order.
		plan[len(lst)-i-1] = &Action{
			Action:    ActionTypeUndo,
			Migration: m,
		}
	}

	return plan, nil
}
