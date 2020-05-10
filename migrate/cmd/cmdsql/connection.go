package cmdsql

import (
	dbSQL "database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/setare/migrations"
	migrationsSQL "github.com/setare/migrations/sql"
)

var (
	connection *dbSQL.DB
)

func Connect(driver, dsn string) error {
	// Starts the connection
	db, err := dbSQL.Open(driver, dsn)
	if err != nil {
		return err
	}

	// Tests the connection
	err = db.Ping()
	if err != nil {
		return err
	}

	// Save the db ref for later use.
	connection = db
	return nil
}

func Initialize(dir string) (migrations.Source, migrations.Target, error) {

	// Initialize source
	s, err := migrations.NewSourceSQLFromDir(dir)
	if err != nil {
		return nil, nil, err
	}

	// Initialize target
	t, err := migrationsSQL.NewTarget(s, connection)
	if err != nil {
		return nil, nil, err
	}

	err = t.Create()
	if err != nil {
		return nil, nil, err
	}

	return s, t, nil
}

func NewExecutionContext() migrations.ExecutionContext {
	return migrations.NewExecutionContext(connection)
}
