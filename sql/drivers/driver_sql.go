package drivers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jamillosantos/migrations/v2"

	"github.com/spaolacci/murmur3"
)

var (
	ErrMissingDatabaseName     = errors.New("missing database name")
	ErrFailedToGetAffectedRows = errors.New("failed getting the number of affected rows")
)

type TXExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

type sqlDriver struct {
	db           DB
	databaseName string
	tableName    string
}

func newSQL(db DB, options ...Option) (Driver, error) {
	opts := driverOpts{
		TableName: DefaultMigrationsTableName,
	}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.DatabaseName == "" {
		return nil, ErrMissingDatabaseName
	}

	return &sqlDriver{
			db:           db,
			databaseName: opts.DatabaseName,
			tableName:    opts.TableName,
		},
		nil
}

// noopUnlocker is a dummy implementation of the migrations.Unlocker interface that does nothing. This is used for the
// sqlDriver.Lock method as the locker mechanism is very database specific.
type noopUnlocker struct {
}

func (n noopUnlocker) Unlock(_ context.Context) error {
	return nil
}

func (p *sqlDriver) Lock(_ context.Context) (migrations.Unlocker, error) {
	return &noopUnlocker{}, nil
}

func (p *sqlDriver) Add(ctx context.Context, id string) error {
	_, err := p.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (id, dirty) VALUES ($1, true)", p.tableName), id)
	if err != nil {
		return fmt.Errorf("failed adding migration to the executed list: %w", err)
	}
	return nil
}

func (p *sqlDriver) Remove(ctx context.Context, id string) error {
	result, err := p.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", p.tableName), id)
	if err != nil {
		return fmt.Errorf("failed removing migration from the executed list: %w", err)
	}

	if rows, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetAffectedRows, err)
	} else if rows == 0 {
		return migrations.ErrMigrationNotFound
	}
	return nil
}

func (p *sqlDriver) StartMigration(ctx context.Context, id string) error {
	result, err := p.db.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET dirty = true WHERE id = $1", p.tableName), id)
	if err != nil {
		return fmt.Errorf("failed starting migration: %w", err)
	}
	if rows, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetAffectedRows, err)
	} else if rows == 0 {
		return migrations.ErrMigrationNotFound
	}
	return nil
}

func (p *sqlDriver) FinishMigration(ctx context.Context, id string) error {
	result, err := p.db.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET dirty = false WHERE id = $1", p.tableName), id)
	if err != nil {
		return fmt.Errorf("failed finishing migration: %w", err)
	}
	if rows, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetAffectedRows, err)
	} else if rows == 0 {
		return migrations.ErrMigrationNotFound
	}
	return nil
}

func (p *sqlDriver) generateLockID() (int64, error) {
	h := murmur3.New64()
	if _, err := h.Write([]byte(p.databaseName)); err != nil {
		return 0, err
	}
	if _, err := h.Write([]byte("|||")); err != nil {
		return 0, err
	}
	if _, err := h.Write([]byte(p.tableName)); err != nil {
		return 0, err
	}
	return int64(h.Sum64()), nil
}
