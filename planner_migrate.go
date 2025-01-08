package migrations

import (
	"context"
	"errors"
	"fmt"
)

type migratePlanner struct {
	source Source
	target Target
}

// MigratePlanner is an ActionPlanner that returns a Planner that plans actions to take the current version of the
// database to the latest.
func MigratePlanner(source Source, target Target) Planner {
	return &migratePlanner{
		source: source,
		target: target,
	}
}

func (planner *migratePlanner) Plan(ctx context.Context) (Plan, error) {
	repo, err := planner.source.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error listing available migrations", err)
	}

	migrationList, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}

	currentMigrationID, err := planner.target.Current(ctx)
	if errors.Is(err, ErrNoCurrentMigration) {
		plan := make(Plan, len(migrationList))
		for i, m := range migrationList {
			plan[i] = &Action{
				Action:    ActionTypeDo,
				Migration: m,
			}
		}
		// If there is no current migration, all migrations should run.
		return plan, nil
	} else if err != nil {
		// Otherwise, it is just an error.
		return nil, err
	}

	done, err := planner.target.Done(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing migrations applied: %w", err)
	}
	for _, migrationID := range done {
		_, err := repo.ByID(migrationID)
		if err != nil {
			return nil, err
		}
	}

	// Detects if a migration was added to the list of applied migrations that is not in the source list.
	for i, m := range done {
		if m != migrationList[i].ID() {
			return nil, fmt.Errorf("%w: %s", ErrStaleMigrationDetected, migrationList[i].String())
		}
	}

	// This is the migration that we are trying to reach. Always the most recent one.
	targetMigration := migrationList[len(migrationList)-1]

	// If the current migration is the same as the target migration
	if currentMigrationID == targetMigration.ID() {
		// Nothing should be done.
		return Plan{}, nil
	}

	currentMigrationIndex, err := findMigrationIndexByID(migrationList, currentMigrationID)
	if err != nil {
		return nil, err
	}

	// If the current migration is further in the future than the target migration.
	if currentMigrationID > targetMigration.ID() {
		return nil, fmt.Errorf("%w: current %s, target %s", ErrCurrentMigrationMoreRecent, currentMigrationID, targetMigration.ID())
	}

	// Build plan
	lst := migrationList[currentMigrationIndex+1:]
	plan := make(Plan, len(lst))
	for i, m := range lst {
		plan[i] = &Action{
			Action:    ActionTypeDo,
			Migration: m,
		}
	}

	return plan, nil
}
