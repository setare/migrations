package migrations

import "time"

// Source lists all migrations.
//
// Migrations can be stored into many medias, from Go Source files until plain
// SQL files. This interface is responsible for abstracting how this system
// accepts any media to list the migrations.
type Source interface {
	// ByID will return the Migration reference given the ID.
	ByID(time.Time) (Migration, error)

	// List will returns the list of available migrations.
	//
	// If there is no migrations available, an `ErrNoMigrationsAvailable` should
	// be returned.
	List() ([]Migration, error)
}
