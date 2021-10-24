package migrations

import (
	"errors"
)

type rewindPlanner struct {
	source Source
	target Target
}

func RewindPlanner(source Source, target Target) Planner {
	return &rewindPlanner{
		source: source,
		target: target,
	}
}

func (planner *rewindPlanner) Plan() (Plan, error) {
	list, err := planner.source.List()
	if err != nil {
		return nil, err
	}

	current, err := planner.target.Current()
	if errors.Is(err, ErrNoCurrentMigration) {
		// If there is no current migration, no migrations should run.
		return Plan{}, nil
	} else if err != nil {
		// Otherwise, it is just an error.
		return nil, err
	}

	currentMigrationIndex, err := findMigrationIndex(list, current)
	if err != nil {
		return nil, err
	}

	// Build the plan
	lst := list[:currentMigrationIndex+1]
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
