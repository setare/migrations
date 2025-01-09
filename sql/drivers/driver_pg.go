package drivers

import (
	"context"
	"fmt"

	"github.com/jamillosantos/migrations/v2"
)

type pgDriver struct {
	sqlDriver
}

func newPostgres(db DB, options ...Option) (Driver, error) {
	opts := driverOpts{
		TableName: DefaultMigrationsTableName,
	}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	if opts.DatabaseName == "" {
		rows, err := db.QueryContext(opts.Ctx, "SELECT current_database()")
		if err != nil {
			return nil, fmt.Errorf("error obtaining current database name: %w", err)
		}
		defer func() {
			_ = rows.Close()
		}()

		if rows.Next() {
			err = rows.Scan(&opts.DatabaseName)
			if err != nil {
				return nil, fmt.Errorf("error scanning current database name: %w", err)
			}
		} else {
			return nil, ErrMissingDatabaseName
		}
	}

	return &pgDriver{
		sqlDriver{
			db:           db,
			databaseName: opts.DatabaseName,
			tableName:    opts.TableName,
		},
	}, nil
}

// pgLocker is the migrations.Locker implementation for sqlDriver database. Its job is to block other instances of the
// migration system to run at the same time. In other to achieve this, it uses the database and table name to create a
// unique key that is hashed (using murmur3) to a bigint. Then, an advisory lock is created using that key.
type pgLocker struct {
	db   TXExecer
	code int64
}

func (p *pgLocker) Unlock(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", p.code)
	if err != nil {
		_ = p.db.Rollback()
		return fmt.Errorf("failed unlocking migration: %w", err)
	}
	_ = p.db.Commit()
	return nil
}

func (p *pgDriver) Lock(ctx context.Context) (migrations.Unlocker, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction for locking: %w", err)
	}

	advisoryLockID, err := p.generateLockID()
	if err != nil {
		return nil, fmt.Errorf("failed locking database: %w", err)
	}
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_lock($1)", advisoryLockID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed locking database: %w", err)
	}
	return &pgLocker{db: tx, code: advisoryLockID}, nil
}
