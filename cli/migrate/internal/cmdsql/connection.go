package cmdsql

import (
	dbSQL "database/sql"
	"errors"

	"github.com/briandowns/spinner"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jamillosantos/migrations"
	"github.com/jamillosantos/migrations/cmd/uiutils"
	migrationsSQL "github.com/jamillosantos/migrations/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ory/viper"
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

func Initialize() (migrations.Source, migrations.Target, error) {
	dir := viper.GetString("directory")

	if !viper.IsSet("dsn") || viper.Get("dsn") == "" {
		return nil, nil, errors.New("--dsn or DSN environment variable not defined")
	}

	err := uiutils.Spin(func(s *spinner.Spinner) error {
		s.Suffix = "Connecting ..."
		return Connect(viper.GetString("driver"), viper.GetString("dsn"))
	})
	if err != nil {
		return nil, nil, err
	}

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
