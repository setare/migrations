package migrations

import (
	"errors"
)

type resetPlanner struct {
	source Source
	target Target
}

func ResetPlanner(source Source, target Target) Planner {
	return &resetPlanner{
		source: source,
		target: target,
	}
}

func (planner *resetPlanner) Plan() (Plan, error) {
	list, err := planner.source.List()
	if err != nil {
		return nil, err
	}

	current, err := planner.target.Current()
	if errors.Is(err, ErrNoCurrentMigration) {
		plan := make(Plan, len(list))
		for i, m := range list {
			plan[i] = &Action{
				Action:    ActionTypeDo,
				Migration: m,
			}
		}
		// If there is no current migration, no migrations should run.
		return plan, nil
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
	plan := make(Plan, len(lst), len(lst)+len(list))
	// Rewinds it...
	for i, m := range lst {
		// Inverts the order of the list, the rewind should be planned in the inverse execution order.
		plan[len(lst)-i-1] = &Action{
			Action:    ActionTypeUndo,
			Migration: m,
		}
	}
	// Migrates it ...
	for _, m := range list {
		plan = append(plan, &Action{
			Action:    ActionTypeDo,
			Migration: m,
		})
	}

	return plan, nil
}
