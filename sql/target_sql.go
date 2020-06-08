package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/setare/migrations"
)

type targetSQL struct {
	source    migrations.Source
	db        *sql.DB
	tableName string
}

// Option
type Option func(target *targetSQL) error

func NewTarget(source migrations.Source, db *sql.DB, options ...Option) (migrations.Target, error) {
	target := &targetSQL{
		source,
		db,
		"_migrations",
	}
	for _, opt := range options {
		err := opt(target)
		if err != nil {
			return nil, err
		}
	}
	return target, nil
}

func Table(tableName string) Option {
	return func(target *targetSQL) error {
		target.tableName = tableName
		return nil
	}
}

// OptError is an Option that will return the given error when initializing the
// target. That is really useful for testing.
func OptError(err error) Option {
	return func(target *targetSQL) error {
		return err
	}
}

func (target *targetSQL) Create() error {
	_, err := target.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id BIGINT PRIMARY KEY)", target.tableName))
	return err
}

func (target *targetSQL) Destroy() error {
	_, err := target.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", target.tableName))
	return err
}

func (target *targetSQL) Current() (migrations.Migration, error) {
	list, err := target.Done()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, migrations.ErrNoCurrentMigration
	}
	return list[len(list)-1], nil
}

func (target *targetSQL) Done() ([]migrations.Migration, error) {
	rs, err := target.db.Query(fmt.Sprintf("SELECT id FROM %s ORDER BY id ASC", target.tableName))
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var id int64
	result := make([]migrations.Migration, 0)
	for rs.Next() {
		err := rs.Scan(&id)
		if err != nil {
			return nil, err
		}
		idDt := time.Unix(id, 0)
		migration, err := target.source.ByID(idDt)
		if err != nil {
			return nil, err
		}
		result = append(result, migration)
	}
	return result, nil
}

func (target *targetSQL) Add(migration migrations.Migration) error {
	_, err := target.db.Exec(fmt.Sprintf("INSERT INTO %s (id) VALUES (%d)", target.tableName, migration.ID().Unix()))
	return err
}

func (target *targetSQL) Remove(migration migrations.Migration) error {
	_, err := target.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = %d", target.tableName, migration.ID().UTC().Unix()))
	return err
}
