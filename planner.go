package migrations

import "github.com/pkg/errors"

type Planner interface {
	Plan() (Plan, error)
}

type Plan []*Action

type PlannerFunc func(Source, Target) Planner

type basePlanner struct {
	source Source
	target Target
}

func (planner *basePlanner) getMigrationIndex(list []Migration, migration Migration) (int, error) {
	index := -1
	for i, m := range list {
		if migration.ID() == m.ID() {
			index = i
		}
	}

	// If the current migration is not in the list
	if index == -1 {
		return -1, errors.Wrap(ErrCurrentMigrationNotFound, migration.ID().Format(DefaultMigrationIDFormat))
	}

	return index, nil
}

type migratePlanner struct {
	basePlanner
}

func MigratePlanner(source Source, target Target) Planner {
	return &migratePlanner{
		basePlanner{
			source: source,
			target: target,
		},
	}
}

func (planner *migratePlanner) Plan() (Plan, error) {
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
		// If there is no current migration, all migrations should run.
		return plan, nil
	} else if err != nil {
		// Otherwise, it is just an error.
		return nil, err
	}

	// This is the migration that we are trying to reach.
	targetMigration := list[len(list)-1]

	// If the current migration is the same as the target migration
	if current.ID() == targetMigration.ID() {
		// Nothing should be done.
		return Plan{}, nil
	}

	currentMigrationIndex, err := planner.getMigrationIndex(list, current)
	if err != nil {
		return nil, err
	}

	// If the current migration is further in the future than the target migration.
	if current.ID().After(targetMigration.ID()) {
		return nil, errors.Wrapf(ErrCurrentMigrationMoreRecent, "current %s, target %s", current.ID().Format(DefaultMigrationIDFormat), targetMigration.ID().Format(DefaultMigrationIDFormat))
	}

	// Build plan
	lst := list[currentMigrationIndex+1:]
	plan := make(Plan, len(lst))
	for i, m := range lst {
		plan[i] = &Action{
			Action:    ActionTypeDo,
			Migration: m,
		}
	}

	return plan, nil
}

type rewindPlanner struct {
	basePlanner
}

func RewindPlanner(source Source, target Target) Planner {
	return &rewindPlanner{
		basePlanner{
			source: source,
			target: target,
		},
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

	currentMigrationIndex, err := planner.getMigrationIndex(list, current)
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

type resetPlanner struct {
	basePlanner
}

func ResetPlanner(source Source, target Target) Planner {
	return &resetPlanner{
		basePlanner{
			source: source,
			target: target,
		},
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

	currentMigrationIndex, err := planner.getMigrationIndex(list, current)
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

type stepPlanner struct {
	basePlanner
	step int
}

func StepPlanner(step int) PlannerFunc {
	return func(source Source, target Target) Planner {
		return &stepPlanner{
			basePlanner: basePlanner{
				source: source,
				target: target,
			},
			step: step,
		}
	}
}

func DoPlanner(source Source, target Target) Planner {
	return StepPlanner(1)(source, target)
}

func UndoPlanner(source Source, target Target) Planner {
	return StepPlanner(-1)(source, target)
}

func (planner *stepPlanner) Plan() (Plan, error) {
	list, err := planner.source.List()
	if err != nil {
		return nil, err
	}

	current, err := planner.target.Current()
	if err != nil && ((planner.step < 0 && err == ErrNoCurrentMigration) || err != ErrNoCurrentMigration) {
		return nil, err
	}

	plan := make(Plan, 0)
	if planner.step < 0 {
		currentMigrationIndex, err := planner.getMigrationIndex(list, current)
		if err != nil {
			return nil, err
		}

		if currentMigrationIndex+planner.step+1 < 0 {
			return nil, ErrStepOutOfIndex
		}

		for i := currentMigrationIndex; i > currentMigrationIndex+planner.step; i-- {
			m := list[i]
			plan = append(plan, &Action{
				Action:    ActionTypeUndo,
				Migration: m,
			})
		}
	} else if planner.step > 0 {
		var currentMigrationIndex int = -1
		if current != nil {
			currentMigrationIndex, err = planner.getMigrationIndex(list, current)
			if err != nil {
				return nil, err
			}
		}

		if currentMigrationIndex+planner.step >= len(list) {
			return nil, ErrStepOutOfIndex
		}

		if current == nil {
			lst := list[:planner.step]
			plan := make(Plan, planner.step)
			for i, m := range lst {
				plan[i] = &Action{
					Action:    ActionTypeDo,
					Migration: m,
				}
			}
			return plan, nil
		}

		for i := currentMigrationIndex; i < currentMigrationIndex+planner.step; i++ {
			plan = append(plan, &Action{
				Action:    ActionTypeDo,
				Migration: list[i],
			})
		}
	}
	return plan, nil
}
