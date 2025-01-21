package migrations

import (
	"context"
	"fmt"
)

type migrationFunc func(ctx context.Context) error

type BaseMigration struct {
	id          string
	description string
	do          migrationFunc
	undo        migrationFunc
	next        Migration
	previous    Migration
}

func NewMigration(id, description string, do, undo migrationFunc) *BaseMigration {
	return &BaseMigration{
		id:          id,
		description: description,
		do:          do,
		undo:        undo,
	}
}

// ID identifies the migration. Through the ID, all the sorting is done.
func (migration *BaseMigration) ID() string {
	return migration.id
}

// String will return a representation of the migration into a string format
// for user identification.
func (migration *BaseMigration) String() string {
	return fmt.Sprintf("%s_%s", migration.id, migration.description)
}

// Description is the humanized description for the migration.
func (migration *BaseMigration) Description() string {
	return migration.description
}

// Next will link this migration with the next. This link should be created
// by the source while it is being loaded.
func (migration *BaseMigration) Next() Migration {
	return migration.next
}

// SetNext will set the next migration
func (migration *BaseMigration) SetNext(value Migration) Migration {
	migration.next = value
	return migration
}

// Previous will link this migration with the previous. This link should be
// created by the Source while it is being loaded.
func (migration *BaseMigration) Previous() Migration {
	return migration.previous
}

// SetPrevious will set the previous migration
func (migration *BaseMigration) SetPrevious(value Migration) Migration {
	migration.previous = value
	return migration
}

// Do will execute the migration.
func (migration *BaseMigration) Do(ctx context.Context) error {
	return migration.do(ctx)
}

// CanUndo is a flag that mark this flag as undoable.
func (migration *BaseMigration) CanUndo() bool {
	return migration.undo != nil
}

// Undo will undo the migration.
func (migration *BaseMigration) Undo(ctx context.Context) error {
	return migration.undo(ctx)
}
