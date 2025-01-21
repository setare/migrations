package sql

import (
	"context"
	"fmt"

	"github.com/jamillosantos/migrations/v2"
)

type migrationSQL struct {
	dbGetter func() DBExecer

	id              string
	description     string
	next            migrations.Migration
	previous        migrations.Migration
	doFile          string
	doFileContent   string
	undoFile        string
	undoFileContent string
}

// ID identifies the migration. Through the ID, all the sorting is done.
func (migration *migrationSQL) ID() string {
	return migration.id
}

// String will return a representation of the migration into a string format
// for user identification.
func (migration *migrationSQL) String() string {
	if migration.CanUndo() {
		return fmt.Sprintf("[%s,%s]", migration.doFile, migration.undoFile)
	}
	return fmt.Sprintf("[%s]", migration.doFile)
}

// Description is the humanized description for the migration.
func (migration *migrationSQL) Description() string {
	return migration.description
}

// Next will link this migration with the next. This link should be created
// by the source while it is being loaded.
func (migration *migrationSQL) Next() migrations.Migration {
	return migration.next
}

// SetNext will set the next migration
func (migration *migrationSQL) SetNext(value migrations.Migration) migrations.Migration {
	migration.next = value
	return migration
}

// Previous will link this migration with the previous. This link should be
// created by the Source while it is being loaded.
func (migration *migrationSQL) Previous() migrations.Migration {
	return migration.previous
}

// SetPrevious will set the previous migration
func (migration *migrationSQL) SetPrevious(value migrations.Migration) migrations.Migration {
	migration.previous = value
	return migration
}

func (migration *migrationSQL) executeSQL(ctx context.Context, sql string) error {
	db := migration.dbGetter()

	_, err := db.ExecContext(ctx, sql)
	if err != nil {
		return migrations.NewQueryError(err, sql)
	}
	return nil
}

// Do will execute the migration.
func (migration *migrationSQL) Do(ctx context.Context) error {
	return migration.executeSQL(ctx, migration.doFileContent)
}

// CanUndo is a flag that mark this flag as undoable.
func (migration *migrationSQL) CanUndo() bool {
	return migration.undoFile != ""
}

// Undo will undo the migration.
func (migration *migrationSQL) Undo(ctx context.Context) error {
	return migration.executeSQL(ctx, migration.undoFileContent)
}
