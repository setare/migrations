package migrations

import (
	"errors"
	"strings"
)

var (
	ErrMigrationNotFound          = errors.New("migration not found")
	ErrNoCurrentMigration         = errors.New("no current migration")
	ErrCurrentMigrationNotFound   = errors.New("current migration not found in the list")
	ErrCurrentMigrationMoreRecent = errors.New("current migration is more recent than target migration")
	ErrNoMigrationsAvailable      = errors.New("no migrations available")
	ErrMigrationNotUndoable       = errors.New("migration cannot be undone")

	// ErrStepOutOfIndex is returned when a `StepResolver` cannot resolve a
	// migration due to the resolved index be outside of the migration list.
	ErrStepOutOfIndex = errors.New("step out of bounds")

	// ErrMigrationNotListed is returned when a migration is not found in the
	// `Source` list.
	ErrMigrationNotListed = errors.New("migration not in the source list")

	// ErrStaleMigrationDetected is returned when a migration with an ID eariler of the current applied migration is
	// detected.
	ErrStaleMigrationDetected = errors.New("stale migration detected")

	// ErrInvalidAction is returned when, while executing, the `Action.Action`
	// has an invalid value.
	ErrInvalidAction = errors.New("undefined action")
)

// ---------------------------------------------------------------------------------------------------------------------

// MigrationIDError wraps an error with a migration ID.
type MigrationIDError interface {
	MigrationID() string
}

// MigrationError wraps an error with a migration property.
type MigrationError interface {
	Migration() Migration
	Unwrap() error
}

// MigrationsError wraps an error with a list of migrations.
type MigrationsError interface {
	Migrations() []Migration
}

// ---------------------------------------------------------------------------------------------------------------------

type migrationError struct {
	error
	migration Migration
}

// WrapMigration creates a `MigrationError` based on an existing error.
func WrapMigration(err error, migration Migration) *migrationError {
	return &migrationError{
		err,
		migration,
	}
}

func (err *migrationError) Migration() Migration {
	return err.migration
}

func (err *migrationError) Unwrap() error {
	return err.error
}

func (err *migrationError) Error() string {
	return err.migration.ID() + ": " + err.error.Error()
}

func (err *migrationError) Is(target error) bool {
	return err == target || errors.Is(err.error, target)
}

// ---------------------------------------------------------------------------------------------------------------------

type migrationsError struct {
	error
	migrations []Migration
}

func WrapMigrations(err error, migrations ...Migration) MigrationsError {
	return &migrationsError{
		err,
		migrations,
	}
}

func (err *migrationsError) Migrations() []Migration {
	return err.migrations
}

func (err *migrationsError) Error() string {
	var r strings.Builder
	for i, migration := range err.migrations {
		if i > 0 {
			r.WriteString(",")
		}
		r.WriteString(migration.ID())
	}
	r.WriteString(": ")
	r.WriteString(err.error.Error())
	return r.String()
}

func (err *migrationsError) Unwrap() error {
	return err.error
}

func (err *migrationsError) Is(target error) bool {
	return err == target || errors.Is(err.error, target)
}

// ---------------------------------------------------------------------------------------------------------------------

type migrationIDError struct {
	error
	migrationID string
}

// WrapMigrationID creates a `MigrationCodeError` based on an existing error.
func WrapMigrationID(err error, migrationID string) *migrationIDError {
	return &migrationIDError{
		err,
		migrationID,
	}
}

func (err *migrationIDError) MigrationID() string {
	return err.migrationID
}

func (err *migrationIDError) Unwrap() error {
	return err.error
}

func (err *migrationIDError) Is(target error) bool {
	if target, ok := target.(MigrationIDError); ok {
		return target.MigrationID() == err.migrationID
	}
	return errors.Is(err.error, target)
}

func (err *migrationIDError) Error() string {
	return "migration " + err.migrationID + ": " + err.error.Error()
}

// ---------------------------------------------------------------------------------------------------------------------

type QueryError interface {
	Query() string
}

type queryError struct {
	error
	query string
}

func NewQueryError(err error, query string) error {
	return &queryError{err, query}
}

func (err *queryError) Unwrap() error {
	return err.error
}

func (err *queryError) Query() string {
	return err.query
}

func (err *queryError) Is(target error) bool {
	return err == target || errors.Is(err.error, target)
}
