package migrations

import "github.com/pkg/errors"

// MigrationResolver is used by the `Planner` to retrieve the target migrations.
//
// The `Planner` will receive an array of `MigrationResolver`s and will use its
// resolutions to go back and forth in the history.
type MigrationResolver interface {
	// Resolve returns the migrations that the `Planner` should try to reach.
	// It can be prior or posterior the current migration.
	Resolve() (Migration, error)
}

type baseResolver struct {
	source Source
	target Target
}

func newResolver(source Source, target Target) *baseResolver {
	return &baseResolver{
		source: source,
		target: target,
	}
}

type resolverStep struct {
	*baseResolver
	stepSize int
}

// StepResolver will return a `MigrationResolver` that takes into consideration
// the current migration index and make it walk forward X steps. The number of
// steps can be positive (going forward) or negative (undoing migrations).
func StepResolver(source Source, target Target, stepSize int) MigrationResolver {
	return &resolverStep{
		baseResolver: newResolver(source, target),
		stepSize:     stepSize,
	}
}

// Resolve will get the current migration index and calculates what is the
// migration given how many steps should be performed. The number of steps can
// be positive (going forward) or negative (undoing migrations).
//
// This function will not take into consideration undoable migrations.
func (resolver *resolverStep) Resolve() (Migration, error) {
	currentMigration, err := resolver.target.Current()
	if err != ErrNoCurrentMigration && err != nil {
		return nil, err
	}

	list, err := resolver.source.List()
	if err != nil {
		return nil, err
	}

	migrationIdx := -1
	if currentMigration != nil {
		migrationIdx = getMigrationIdx(currentMigration, list)
		if migrationIdx == -1 {
			return nil, WrapMigration(ErrMigrationNotFound, currentMigration)
		}
	}

	targetIdx := migrationIdx + resolver.stepSize

	if targetIdx < 0 || targetIdx >= len(list) {
		s := "+"
		if resolver.stepSize < 0 {
			s = "-"
		}
		return nil, errors.Wrapf(ErrStepOutOfIndex, "%s%d cannot be resolved", s, resolver.stepSize)
	}

	return list[targetIdx], nil
}

type mostRecentResolver struct {
	*baseResolver
}

func MostRecentResolver(source Source) MigrationResolver {
	return &mostRecentResolver{
		baseResolver: newResolver(source, nil),
	}
}

func (resolver *mostRecentResolver) Resolve() (Migration, error) {
	list, err := resolver.source.List()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNoMigrationsAvailable
	}
	return list[len(list)-1], nil
}

type firstMigrationResolver struct {
	*baseResolver
}

func FirstMigrationResolver(source Source) MigrationResolver {
	return &firstMigrationResolver{
		baseResolver: newResolver(source, nil),
	}
}

func (resolver *firstMigrationResolver) Resolve() (Migration, error) {
	list, err := resolver.source.List()
	if err != nil {
		return nil, err
	}
	return list[0], nil
}
