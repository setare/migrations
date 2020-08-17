package testingutils

import (
	"fmt"
	"time"

	"github.com/jamillosantos/migrations"
)

type DumpForwardMigration struct {
	id          time.Time
	DoneCount   int
	DoErr       error
	UndoneCount int
	UndoErr     error
	previous    migrations.Migration
	next        migrations.Migration
}

type DumpMigration struct {
	DumpForwardMigration
}

func NewForwardMigration(id time.Time, err ...error) *DumpForwardMigration {
	var doErr error
	if len(err) > 0 {
		doErr = err[0]
	}
	return &DumpForwardMigration{
		id:    id,
		DoErr: doErr,
	}
}

func NewMigration(id time.Time, err ...error) *DumpMigration {
	var doErr, undoErr error
	if len(err) > 0 {
		doErr = err[0]
	}
	if len(err) > 1 {
		undoErr = err[0]
	}
	return &DumpMigration{
		DumpForwardMigration: DumpForwardMigration{
			id:      id,
			DoErr:   doErr,
			UndoErr: undoErr,
		},
	}
}

func (migration *DumpForwardMigration) ID() time.Time {
	return migration.id
}

func (migration *DumpForwardMigration) String() string {
	return fmt.Sprintf("[%s]", migration.id.Format(migrations.DefaultMigrationIDFormat))
}

func (migration *DumpForwardMigration) Description() string {
	return "dumb"
}

func (migration *DumpForwardMigration) Do(migrations.ExecutionContext) error {
	migration.DoneCount++
	return migration.DoErr
}

func (migration *DumpMigration) CanUndo() bool {

	return true
}

func (migration *DumpForwardMigration) Undo(migrations.ExecutionContext) error {
	migration.UndoneCount++
	return migration.UndoErr
}

func (migration *DumpForwardMigration) CanUndo() bool {
	return false
}

func (migration *DumpForwardMigration) Next() migrations.Migration {
	return migration.next
}

func (migration *DumpForwardMigration) SetNext(next migrations.Migration) migrations.Migration {
	migration.next = next
	return migration
}

func (migration *DumpForwardMigration) Previous() migrations.Migration {
	return migration.previous
}

func (migration *DumpForwardMigration) SetPrevious(previous migrations.Migration) migrations.Migration {
	migration.previous = previous
	return migration
}
