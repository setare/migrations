package migrations

import (
	"context"
	"errors"
)

type resetPlanner struct {
	source Source
	target Target
}

// ResetPlanner is an ActionPlanner that returns a Planner that plans actions to rewind all migrations and then migrate
// again to the latest version.
func ResetPlanner(source Source, target Target) Planner {
	return &resetPlanner{
		source: source,
		target: target,
	}
}

func (planner *resetPlanner) Plan(ctx context.Context) (Plan, error) {
	repo, err := planner.source.Load(ctx)
	if err != nil {
		return nil, err
	}

	migrationList, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}

	currentMigrationID, err := planner.target.Current(nil)
	if errors.Is(err, ErrNoCurrentMigration) {
		plan := make(Plan, len(migrationList))
		for i, m := range migrationList {
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

	currentMigrationIndex, err := findMigrationIndexByID(migrationList, currentMigrationID)
	if err != nil {
		return nil, err
	}

	// Build the plan
	lst := migrationList[:currentMigrationIndex+1]
	plan := make(Plan, len(lst), len(lst)+len(migrationList))
	// Rewinds it...
	for i, m := range lst {
		// Inverts the order of the list, the rewind should be planned in the inverse execution order.
		plan[len(lst)-i-1] = &Action{
			Action:    ActionTypeUndo,
			Migration: m,
		}
	}
	// Migrates it ...
	for _, m := range migrationList {
		plan = append(plan, &Action{
			Action:    ActionTypeDo,
			Migration: m,
		})
	}

	return plan, nil
}
