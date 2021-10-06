package migrations

// ListSource is the default Source implementation that keep all the migrations on a list.
//
// This source is meant to be used as base implementation for other specific Sources. Check out the
// sql.Source implementation.
type ListSource struct {
	list []Migration
	byID map[string]Migration
}

// NewSource returns a new instance of the default implementation of the Source.
func NewSource() *ListSource {
	return &ListSource{
		list: make([]Migration, 0),
		byID: make(map[string]Migration, 0),
	}
}

// DefaultSource is the default source instance.
var DefaultSource = NewSource()

// Add adds a Migration to this Source. It will
func (source *ListSource) Add(migration Migration) error {
	for i, m := range source.list {
		if m.ID() == migration.ID() { // If there is a migration with the same ID already added
			return WrapMigration(ErrNonUniqueMigrationID, migration)
		}

		if migration.ID() < m.ID() {
			// If the given migration has happened before the entry from the list we need to insert the given migration
			// before of the current entry.

			// Update the migration helpers
			migration.SetPrevious(m.Previous())
			if prev := migration.Previous(); prev != nil {
				prev.SetNext(migration)
			}
			migration.SetNext(m)
			m.SetPrevious(migration)

			// Add migration in the middle of the list
			source.list = append(source.list[:i], append([]Migration{migration}, source.list[i:]...)...)

			// Caches the migration ID.
			source.byID[migration.ID()] = migration

			return nil
		}
	}

	// Reaching here means that the migration will be added in the end of the list.
	if len(source.list) > 0 {
		// If there is something on the list, we need to update the helpers.
		migration.SetPrevious(source.list[len(source.list)-1])
		migration.Previous().SetNext(migration)
	} else {
		migration.SetPrevious(nil)
	}

	// Since migration is now the most recent, there is no "next".
	migration.SetNext(nil)
	source.list = append(source.list, migration)

	// Caches the migration ID
	source.byID[migration.ID()] = migration
	return nil
}

func (source *ListSource) List() ([]Migration, error) {
	if len(source.list) == 0 {
		return nil, ErrNoMigrationsAvailable
	}
	return source.list, nil
}

func (source *ListSource) ByID(id string) (Migration, error) {
	if migration, ok := source.byID[id]; ok {
		return migration, nil
	}
	return nil, WrapMigrationID(ErrMigrationNotFound, id)
}
