package migrations

import (
	"context"
	"errors"
)

type stepPlanner struct {
	source Source
	target Target
	step   int
}

// StepPlanner build an ActionPlanner that will step the migrations by the given number.
// If the given number is positive, then the step will be forward, otherwise, it will be backward.
func StepPlanner(step int) ActionPLanner {
	return func(source Source, target Target) Planner {
		return &stepPlanner{
			source: source,
			target: target,
			step:   step,
		}
	}
}

func DoPlanner(source Source, target Target) Planner {
	return StepPlanner(1)(source, target)
}

func UndoPlanner(source Source, target Target) Planner {
	return StepPlanner(-1)(source, target)
}

func (planner *stepPlanner) Plan(ctx context.Context) (Plan, error) {
	repo, err := planner.source.Load(ctx)
	if err != nil {
		return nil, err
	}

	migrationList, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}

	currentMigrationID, err := planner.target.Current(ctx)
	if err != nil && ((planner.step < 0 && errors.Is(err, ErrNoCurrentMigration)) || errors.Is(err, ErrNoCurrentMigration)) {
		return nil, err
	}

	plan := make(Plan, 0)
	if planner.step < 0 {
		currentMigrationIndex, err := findMigrationIndexByID(migrationList, currentMigrationID)
		if err != nil {
			return nil, err
		}

		if currentMigrationIndex+planner.step+1 < 0 {
			return nil, ErrStepOutOfIndex
		}

		for i := currentMigrationIndex; i > currentMigrationIndex+planner.step; i-- {
			m := migrationList[i]
			plan = append(plan, &Action{
				Action:    ActionTypeUndo,
				Migration: m,
			})
		}
	} else if planner.step > 0 {
		var currentMigrationIndex int = -1
		if currentMigrationID != "" {
			currentMigrationIndex, err = findMigrationIndexByID(migrationList, currentMigrationID)
			if err != nil {
				return nil, err
			}
		}

		if currentMigrationIndex+planner.step >= len(migrationList) {
			return nil, ErrStepOutOfIndex
		}

		if currentMigrationID == "" {
			lst := migrationList[:planner.step]
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
				Migration: migrationList[i],
			})
		}
	}
	return plan, nil
}
