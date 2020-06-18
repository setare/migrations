package migrations

import "time"

// Source lists all
//
// Migrations can be stored into many medias, from Go Source files until plain
// SQL files. This interface is responsible for abstracting how this system
// accepts any media to list the
type Source interface {
	// ByID will return the Migration reference given the ID.
	ByID(time.Time) (Migration, error)

	// List will returns the list of available
	//
	// If there is no migrations available, an `ErrNoMigrationsAvailable` should
	// be returned.
	List() ([]Migration, error)

	//
	Add(migration Migration) error
}

type baseSource struct {
	list []Migration
	byID map[int64]Migration
}

func NewSource() Source {
	return &baseSource{
		list: make([]Migration, 0),
		byID: make(map[int64]Migration, 0),
	}
}

var DefaultSource = NewSource()

func (source *baseSource) Add(migration Migration) error {
	for i, m := range source.list {
		if m.ID() == migration.ID() {
			return WrapMigration(ErrNonUniqueMigrationID, migration)
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
			source.list = append(source.list[:i], append([]Migration{migration}, source.list[i:]...)...)
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

func (source *baseSource) List() ([]Migration, error) {
	if len(source.list) == 0 {
		return nil, ErrNoMigrationsAvailable
	}
	return source.list, nil
}

func (source *baseSource) ByID(id time.Time) (Migration, error) {
	if migration, ok := source.byID[id.Unix()]; ok {
		return migration, nil
	}
	return nil, WrapMigrationID(ErrMigrationNotFound, id)
}
