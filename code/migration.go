package code

import (
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/setare/migrations"
)

var (
	ErrInvalidFilename    = errors.New("invalid filename")
	ErrInvalidMigrationID = errors.New("invalid migration ID")
)

type migrationFunc func(migrations.ExecutionContext) error

type codeMigration struct {
	id          time.Time
	description string
	fileName    string
	do          migrationFunc
	undo        migrationFunc
	next        migrations.Migration
	previous    migrations.Migration
}

type Migration struct {
	Skip int
	Do   migrationFunc
	Undo migrationFunc
}

var fnameRegex = regexp.MustCompile("^(\\d+)_(.*)\\.go$")

func MustNew(migration *Migration) migrations.Migration {
	m, err := New(migration)
	if err != nil {
		panic(err)
	}
	return m
}

func New(migration *Migration) (migrations.Migration, error) {
	skip := 2 + migration.Skip

	_, fileName, _, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("error getting file information")
	}

	fileName = path.Base(fileName)

	matches := fnameRegex.FindStringSubmatch(fileName)
	if len(matches) == 0 {
		return nil, errors.Wrap(ErrInvalidFilename, fileName)
	}

	migrationID, err := time.Parse(migrations.DefaultMigrationIDFormat, matches[1])
	if err != nil {
		return nil, errors.Wrap(ErrInvalidMigrationID, matches[1])
	}

	return &codeMigration{
		id:          migrationID,
		fileName:    fileName,
		description: strings.ReplaceAll(matches[2], "_", " "),
		do:          migration.Do,
		undo:        migration.Undo,
	}, nil
}

// ID identifies the migration. Through the ID, all the sorting is done.
func (migration *codeMigration) ID() time.Time {
	return migration.id
}

// String will return a representation of the migration into a string format
// for user identification.
func (migration *codeMigration) String() string {
	return migration.fileName
}

// Description is the humanized description for the migration.
func (migration *codeMigration) Description() string {
	return migration.description
}

// Next will link this migration with the next. This link should be created
// by the source while it is being loaded.
func (migration *codeMigration) Next() migrations.Migration {
	return migration.next
}

// SetNext will set the next migration
func (migration *codeMigration) SetNext(value migrations.Migration) migrations.Migration {
	migration.next = value
	return migration
}

// Previous will link this migration with the previous. This link should be
// created by the Source while it is being loaded.
func (migration *codeMigration) Previous() migrations.Migration {
	return migration.previous
}

// SetPrevious will set the previous migration
func (migration *codeMigration) SetPrevious(value migrations.Migration) migrations.Migration {
	migration.previous = value
	return migration
}

// Do will execute the migration.
func (migration *codeMigration) Do(executionContext migrations.ExecutionContext) error {
	return migration.do(executionContext)
}

// CanUndo is a flag that mark this flag as undoable.
func (migration *codeMigration) CanUndo() bool {
	return migration.undo != nil
}

// Undo will undo the migration.
func (migration *codeMigration) Undo(executionContext migrations.ExecutionContext) error {
	return migration.undo(executionContext)
}
