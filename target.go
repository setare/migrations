package migrations

type Target interface {
	// Current returns the reference to the most recent ran migration.
	//
	// If there is no migration run, the system will return an
	// `ErrNoCurrentMigration` error.
	Current() (Migration, error)

	// Create creates the media for storing the list of all migrations were
	// executed on this target.
	Create() error

	// Destroy removes the list of the migrations that were run.
	Destroy() error

	Done() ([]Migration, error)
	Add(Migration) error
	Remove(Migration) error
}
