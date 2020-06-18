package migrations

import (
	"time"
)

var (
	// DefaultMigrationIDFormat is the default format for the migrations ID.
	DefaultMigrationIDFormat = "20060102150405"
)

// Migration is the abstraction that defines the migration contract.
//
// Migrations, by default, cannot be undone. But, if the final migration
// implement the `MigrationUndoable` interface, the system will let it be
// undone.
type Migration interface {
	// ID identifies the migration. Through the ID, all the sorting is done.
	ID() time.Time

	// String will return a representation of the migration into a string format
	// for user identification.
	String() string

	// Description is the humanized description for the migration.
	Description() string

	// Next will link this migration with the next. This link should be created
	// by the source while it is being loaded.
	Next() Migration

	// SetNext will set the next migration
	SetNext(Migration) Migration

	// Previous will link this migration with the previous. This link should be
	// created by the Source while it is being loaded.
	Previous() Migration

	// SetPrevious will set the previous migration
	SetPrevious(Migration) Migration

	// Do will execute the migration.
	Do(executionContext ExecutionContext) error

	// CanUndo is a flag that mark this flag as undoable.
	CanUndo() bool

	// Undo will undo the migration.
	Undo(executionContext ExecutionContext) error
}
