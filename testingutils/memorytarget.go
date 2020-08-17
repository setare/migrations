package testingutils

import "github.com/jamillosantos/migrations"

type memoryTarget struct {
	done []migrations.Migration
}

func NewMemoryTarget() migrations.Target {
	return &memoryTarget{
		done: make([]migrations.Migration, 0),
	}
}

// Current returns the reference to the most recent ran migration.
//
// If there is no migration run, the system will return an
// `ErrNoCurrentMigration` error.
func (target *memoryTarget) Current() (migrations.Migration, error) {
	if len(target.done) == 0 {
		return nil, migrations.ErrNoCurrentMigration
	}
	return target.done[len(target.done)-1], nil
}

// Create creates the media for storing the list of all migrations were
// executed on this target.
func (target *memoryTarget) Create() error {
	return nil
}

// Destroy removes the list of the migrations that were run.
func (target *memoryTarget) Destroy() error {
	return nil
}

func (target *memoryTarget) Done() ([]migrations.Migration, error) {
	return target.done, nil
}

func (target *memoryTarget) Add(migration migrations.Migration) error {
	target.done = append(target.done, migration)
	return nil
}

func (target *memoryTarget) Remove(migration migrations.Migration) error {
	for i, m := range target.done {
		if m == migration {
			target.done = append(target.done[:i], target.done[i+1:]...)
			return nil
		}
	}
	return migrations.WrapMigration(migrations.ErrMigrationNotFound, migration)
}
