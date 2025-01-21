//go:build integration

package sql_test

import (
	"context"
	stdsql "database/sql"
	"embed"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/jamillosantos/migrations/v2"
	"github.com/jamillosantos/migrations/v2/reporters"
	"github.com/jamillosantos/migrations/v2/sql"
	"github.com/jamillosantos/migrations/v2/sql/drivers"
)

var (
	//go:embed testdata/case1/*.sql
	case1Migrations embed.FS

	//go:embed testdata/case2/*.sql
	case2Migrations embed.FS

	//go:embed testdata/case3/*.sql
	case3Migrations embed.FS
)

type migrationCase string

const (
	migrationCase1 migrationCase = "case1"
	migrationCase2 migrationCase = "case2"
	migrationCase3 migrationCase = "case3"
)

func migrate(t *testing.T, db *stdsql.DB, migrationCase migrationCase) error {
	t.Helper()

	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	var (
		fs         embed.FS
		caseFolder string
	)

	switch migrationCase {
	case migrationCase1:
		fs = case1Migrations
		caseFolder = "case1"
	case migrationCase2:
		fs = case2Migrations
		caseFolder = "case2"
	case migrationCase3:
		fs = case3Migrations
		caseFolder = "case3"
	}

	source, err := sql.SourceFromFS(func() sql.DBExecer {
		return db
	}, fs, "testdata/"+caseFolder)
	require.NoError(t, err)

	target, err := sql.NewTarget(db, sql.WithDriverOptions(drivers.WithDatabaseName("testdb")))
	require.NoError(t, err)

	_, err = migrations.Migrate(ctx, source, target, migrations.WithRunnerOptions(migrations.WithReporter(reporters.NewZapReporter(logger))))
	return err
}

func TestMigrate_Integration(t *testing.T) {
	t.Run("should run case 1", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, migrationCase1)
		require.NoError(t, err)
	})

	t.Run("should run case 1 + case 2", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, migrationCase1)
		require.NoError(t, err)
		err = migrate(t, db, migrationCase2)
		require.NoError(t, err)
	})

	t.Run("should fail running the case 3 after case 1 and case 2", func(t *testing.T) {
		db := createDBConnection(t)
		err := migrate(t, db, migrationCase1)
		require.NoError(t, err)
		err = migrate(t, db, migrationCase2)
		require.NoError(t, err)
		err = migrate(t, db, migrationCase3)
		require.ErrorIs(t, err, migrations.ErrStaleMigrationDetected)
	})
}

func createDBConnection(t *testing.T) *stdsql.DB {
	t.Helper()

	db, err := stdsql.Open("sqlite3", ":memory:")
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
