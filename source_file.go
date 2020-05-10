package migrations

import (
	"database/sql"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type sourceSQLFiles struct {
	list []Migration
	byID map[int64]Migration
}

var (
	fileRegex = regexp.MustCompile("^(\\d+)(.*)\\.(undo|do).sql$")
)

func NewSourceSQLFromDir(dir string) (Source, error) {
	files, err := filepath.Glob(path.Join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}
	return NewSourceSQLFromFiles(files)
}

func NewSourceSQLFromFiles(files []string) (Source, error) {
	source := &sourceSQLFiles{
		list: make([]Migration, 0, len(files)/2+1),
		byID: make(map[int64]Migration, len(files)/2+1),
	}

	for _, file := range files {
		fBase := path.Base(file)
		g := fileRegex.FindStringSubmatch(fBase)
		if len(g) == 0 {
			return nil, errors.Wrap(ErrInvalidPatternForFile, file)
		}

		migrationID, err := time.Parse(DefaultMigrationIDFormat, g[1])
		if err != nil {
			return nil, err
		}

		description := strings.ReplaceAll(g[2][1:], "_", " ")
		actionType := ActionType(g[3])

		m, err := source.ByID(migrationID)
		if errors.Is(err, ErrMigrationNotFound) {
			m = &migrationSQL{
				id:          migrationID,
				description: description,
			}
			err = source.Add(m)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		mSQL := m.(*migrationSQL)

		if actionType == ActionTypeDo && mSQL.doFile == "" {
			mSQL.doFile = file
		} else if actionType == ActionTypeUndo && mSQL.undoFile == "" {
			mSQL.undoFile = file
		} else {
			return nil, errors.Wrap(ErrInvalidSQLFileNameDuplicated, file)
		}
	}
	return source, nil
}

// ByID will return the Migration reference given the ID.
func (source *sourceSQLFiles) ByID(id time.Time) (Migration, error) {
	if migration, ok := source.byID[id.Unix()]; ok {
		return migration, nil
	}
	return nil, WrapMigrationID(ErrMigrationNotFound, id)
}

// List will returns the list of available migrations.
//
// If there is no migrations available, an `ErrNoMigrationsAvailable` should
// be returned.
func (source *sourceSQLFiles) List() ([]Migration, error) {
	return source.list, nil
}

func (source *sourceSQLFiles) Add(migration Migration) error {
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

type migrationSQL struct {
	id          time.Time
	description string
	next        Migration
	previous    Migration
	doFile      string
	undoFile    string
}

// ID identifies the migration. Through the ID, all the sorting is done.
func (migration *migrationSQL) ID() time.Time {
	return migration.id
}

// Description is the humanized description for the migration.
func (migration *migrationSQL) Description() string {
	return migration.description
}

// Next will link this migration with the next. This link should be created
// by the source while it is being loaded.
func (migration *migrationSQL) Next() Migration {
	return migration.next
}

// SetNext will set the next migration
func (migration *migrationSQL) SetNext(value Migration) Migration {
	migration.next = value
	return migration
}

// Previous will link this migration with the previous. This link should be
// created by the Source while it is being loaded.
func (migration *migrationSQL) Previous() Migration {
	return migration.previous
}

// SetPrevious will set the previous migration
func (migration *migrationSQL) SetPrevious(value Migration) Migration {
	migration.previous = value
	return migration
}

func (migration *migrationSQL) executeFile(executionContext ExecutionContext, file string) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	conn, ok := executionContext.Data().(*sql.DB)
	if !ok {
		return ErrNoCurrentMigration
	}

	sql := string(fileContent)

	_, err = conn.Exec(sql)
	if err != nil {
		return NewQueryError(err, sql)
	}
	return nil
}

// Do will execute the migration.
func (migration *migrationSQL) Do(executionContext ExecutionContext) error {
	return migration.executeFile(executionContext, migration.doFile)
}

// CanUndo is a flag that mark this flag as undoable.
func (migration *migrationSQL) CanUndo() bool {
	return migration.undoFile != ""
}

// Undo will undo the migration.
func (migration *migrationSQL) Undo(executionContext ExecutionContext) error {
	return migration.executeFile(executionContext, migration.undoFile)
}
