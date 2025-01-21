package migrations

import (
	"context"
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
	ID() string

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
	Do(ctx context.Context) error

	// CanUndo is a flag that mark this flag as undoable.
	CanUndo() bool

	// Undo will undo the migration.
	Undo(ctx context.Context) error
}

// Source is responsible to list all migrations available to run.
//
// Migrations can be stored into many medias, from Go source code files, plain SQL files, go:embed. So, this interface
// is responsible for abstracting how this system accepts any media to list the
type Source interface {
	// Add will add a Migration to the source list.
	//
	// If the migration id cannot be found, this function should return ErrMigrationAlreadyExists.
	Add(ctx context.Context, migration Migration) error

	// Load will return the list of available migrations, sorted by ID (the older first, the newest last).
	//
	// If there is no migrations available, an ErrNoMigrationsAvailable should
	// be returned.
	Load(ctx context.Context) (Repository, error)
}

// Target is responsible for managing the state of the migration system. This interface abstracts the operations needed
// to list what migrations were successfully executed, what is the current migration, etc.
type Target interface {
	// Current returns the reference to the most recent migration applied to the system.
	//
	// If there is no migration run, the system will return an ErrNoCurrentMigration error.
	Current(ctx context.Context) (string, error)

	// Create ensures the creation of the list of migrations is done successfully. As an example, if this was an SQL
	// database implementation, this method would create the `_migrations` table.
	Create(ctx context.Context) error

	// Destroy removes the list of the applied migrations. As an example, if this was an SQL database implementation,
	// this would drop the `_migrations` table.
	Destroy(ctx context.Context) error

	// Done list all the migrations that were successfully applied.
	Done(ctx context.Context) ([]string, error)

	// Add adds a migration to the list of successful migrations.
	Add(ctx context.Context, id string) error

	// Remove removes a migration from the list of successful migrations.
	Remove(ctx context.Context, id string) error

	// FinishMigration will mark the migration as finished. This is only used when the migration is being added.
	FinishMigration(ctx context.Context, id string) error

	// StartMigration will mark the migration as dirty. This is only used when the migration is being removed.
	StartMigration(ctx context.Context, id string) error

	// Lock will try locking the migration system in such way no other instance of the process can run the migrations.
	Lock(ctx context.Context) (Unlocker, error)
}

// Unlocker abstracts an implementation for unlocking the migration system.
type Unlocker interface {
	Unlock(ctx context.Context) error
}

type ProgressReporter interface {
	SetStep(current int)
	SetSteps(steps []string)
	SetTotal(total int)
	SetProgress(progress int)
}

type ActionType string

const (
	ActionTypeDo   ActionType = "do"
	ActionTypeUndo ActionType = "undo"
)

type Action struct {
	Action    ActionType
	Migration Migration
}
