package code

import (
	"time"

	"github.com/setare/migrations"
)

type SourceCode interface {
	migrations.Source
	Add(migration migrations.Migration) error
}

type sourceCode struct {
	list []migrations.Migration
	byID map[int64]migrations.Migration
}

func NewSource() SourceCode {
	return &sourceCode{
		list: make([]migrations.Migration, 0),
		byID: make(map[int64]migrations.Migration, 0),
	}
}

func (source *sourceCode) Add(migration migrations.Migration) error {
	for i, m := range source.list {
		if m.ID() == migration.ID() {
			return migrations.WrapMigration(migrations.ErrNonUniqueMigrationID, migration)
		}
		if migration.ID().Before(m.ID()) {
			// Update the migration helpers
			migration.SetPrevious(m.Previous())
			if migration.Previous() != nil {
				migration.Previous().SetNext(migration)
			}
			migration.SetNext(m)
			m.SetPrevious(migration)
			//
			source.list = append(source.list[:i], append([]migrations.Migration{migration}, source.list[i:]...)...)
			source.byID[m.ID().Unix()] = m
			return nil
		}
	}
	if len(source.list) > 0 {
		migration.SetPrevious(source.list[len(source.list)-1])
		migration.Previous().SetNext(migration)
	} else {
		migration.SetPrevious(nil)
	}
	migration.SetNext(nil)
	source.list = append(source.list, migration)
	source.byID[migration.ID().Unix()] = migration
	return nil
}

func (source *sourceCode) List() ([]migrations.Migration, error) {
	if len(source.list) == 0 {
		return nil, migrations.ErrNoMigrationsAvailable
	}
	return source.list, nil
}

func (source *sourceCode) ByID(id time.Time) (migrations.Migration, error) {
	if migration, ok := source.byID[id.Unix()]; ok {
		return migration, nil
	}
	return nil, migrations.WrapMigrationID(migrations.ErrMigrationNotFound, id)
}
