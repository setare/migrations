package migrations

import (
	"context"
)

type Planner interface {
	Plan(ctx context.Context) (Plan, error)
}

type Plan []*Action

func findMigrationIndex(list []Migration, migration Migration) (int, error) {
	index := -1
	for i, m := range list {
		if migration.ID() == m.ID() {
			index = i
		}
	}

	// If the current migration is not in the list
	if index == -1 {
		return -1, WrapMigrationID(ErrCurrentMigrationNotFound, migration.ID())
	}

	return index, nil
}

func findMigrationIndexByID(list []Migration, migrationID string) (int, error) {
	index := -1
	for i, m := range list {
		if migrationID == m.ID() {
			index = i
		}
	}

	// If the current migration is not in the list
	if index == -1 {
		return -1, WrapMigrationID(ErrCurrentMigrationNotFound, migrationID)
	}

	return index, nil
}
