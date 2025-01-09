package drivers

import (
	"context"
	"database/sql"

	"github.com/jamillosantos/migrations/v2"
)

type Driver interface {
	Add(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	StartMigration(ctx context.Context, id string) error
	FinishMigration(ctx context.Context, id string) error
	Lock(ctx context.Context) (migrations.Unlocker, error)
}

type DriverConstructor func(db DB, options ...Option) (Driver, error)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error)
}
