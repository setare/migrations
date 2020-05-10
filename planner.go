package migrations

// Planner will create plans (that are array of `Action`) that can be shown to
// the user or executed.
type Planner struct {
	source Source
	target Target
}
type Plan []*Action

func NewPlanner(source Source, target Target) *Planner {
	return &Planner{
		source: source,
		target: target,
	}
}

type PlanRequest struct {
	Resolvers []MigrationResolver
}

func (planner *Planner) Plan(request *PlanRequest) (Plan, error) {
	currentMigration, err := planner.target.Current()
	if err != ErrNoCurrentMigration && err != nil {
		return nil, err
	}

	list, err := planner.source.List()
	if err != nil {
		return nil, err
	}

	idx := getMigrationIdx(currentMigration, list)
	plan := make(Plan, 0)
	for i, resolver := range request.Resolvers {
		migration, err := resolver.Resolve()
		if err != nil {
			return nil, err
		}

		migrationIdx := getMigrationIdx(migration, list)
		if migrationIdx < 0 {
			return nil, WrapMigration(ErrMigrationNotListed, migration)
		}

		if idx == migrationIdx {
			continue
		}

		goingDown := idx > migrationIdx
		increment := 1
		if goingDown {
			increment = -1
		}

		if goingDown && currentMigration != nil {
			plan = append(plan, &Action{
				Migration: currentMigration,
				Action:    ActionTypeUndo,
			})
		}

		for idx != migrationIdx {
			action := &Action{}
			if idx < migrationIdx {
				action.Action = ActionTypeDo
			} else {
				action.Action = ActionTypeUndo
			}
			idx += increment
			action.Migration = list[idx]
			plan = append(plan, action)
		}

		if i < len(request.Resolvers)-1 {
			if goingDown {
				idx--
			} else {
				idx++
			}
		}
	}

	return plan, nil
}

func (planner *Planner) Migrate() (Plan, error) {
	return planner.Plan(&PlanRequest{
		Resolvers: []MigrationResolver{
			MostRecentResolver(planner.source),
		},
	})
}

func (planner *Planner) Rewind() (Plan, error) {
	return planner.Plan(&PlanRequest{
		Resolvers: []MigrationResolver{
			FirstMigrationResolver(planner.source),
		},
	})
}

func (planner *Planner) Reset() (Plan, error) {
	return planner.Plan(&PlanRequest{
		Resolvers: []MigrationResolver{
			FirstMigrationResolver(planner.source),
			MostRecentResolver(planner.source),
		},
	})
}

func (planner *Planner) Step(size int) (Plan, error) {
	return planner.Plan(&PlanRequest{
		Resolvers: []MigrationResolver{
			StepResolver(planner.source, planner.target, size),
		},
	})
}
