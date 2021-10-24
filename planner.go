package migrations

import "github.com/pkg/errors"

type Planner interface {
	Plan() (Plan, error)
}

type Plan []*Action

type PlannerFunc func(Source, Target) Planner

func findMigrationIndex(list []Migration, migration Migration) (int, error) {
	index := -1
	for i, m := range list {
		if migration.ID() == m.ID() {
			index = i
		}
	}

	// If the current migration is not in the list
	if index == -1 {
		return -1, errors.Wrap(ErrCurrentMigrationNotFound, migration.ID())
	}

	return index, nil
}
