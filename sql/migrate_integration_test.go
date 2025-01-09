//go:build integration
// +build integration

package sql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jamillosantos/migrations/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/jamillosantos/migrations/v2/reporters"

	"github.com/jamillosantos/migrations/v2/sql"
)

const (
	base  = "../tests/migrations"
	case1 = base + "/case1"
	case2 = base + "/case2"
	case3 = base + "/case3"
)

func migrate(t *testing.T, db *sql.DB, migrationCase string) error {
	t.Helper()

	dirFS := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), migrationCase))

	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	source, err := sql.FromFS(func() sql.DBExecer {
		return db
	}, dirFS, ".")
	require.NoError(t, err)

	driver, err := sql.NewDriverForSQL(ctx, db, sql.WithDatabaseName("test"))
	require.NoError(t, err)

	target, err := sql.NewTarget(db, driver)
	require.NoError(t, err)

	_, err = migrations.Migrate(ctx, source, target, migrations.WithRunnerOptions(migrations.WithReporter(reporters.NewZapReporter(logger))))
	return err
}

func TestMigrate_Integration(t *testing.T) {
	t.Run("should run case 1", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, case1)
		require.NoError(t, err)
	})

	t.Run("should run case 1 + case 2", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, case1)
		require.NoError(t, err)
		err = migrate(t, db, case2)
		require.NoError(t, err)
	})

	t.Run("should fail running the case 3 after case 1 and case 2", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, case1)
		require.NoError(t, err)
		err = migrate(t, db, case2)
		require.NoError(t, err)
		err = migrate(t, db, case3)
		require.ErrorIs(t, err, migrations.ErrStaleMigrationDetected)
	})
}

func createDBConnection(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
