package migrations

type stepPlanner struct {
	source Source
	target Target
	step   int
}

func StepPlanner(step int) PlannerFunc {
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
		currentMigrationIndex, err := findMigrationIndex(list, current)
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
			currentMigrationIndex, err = findMigrationIndex(list, current)
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
