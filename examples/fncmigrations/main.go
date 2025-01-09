package main

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jamillosantos/migrations/v2"
	fcnmigrations "github.com/jamillosantos/migrations/v2/examples/fncmigrations/migrations"
	"github.com/jamillosantos/migrations/v2/reporters"
	migrationsql "github.com/jamillosantos/migrations/v2/sql"
)

//go:embed migrations/*.sql
var migrationsFolder embed.FS

func main() {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	s, err := migrationsql.SourceFromFS(func() migrationsql.DBExecer {
		return db
	}, migrationsFolder, "migrations")
	if err != nil {
		panic(err)
	}

	for _, m := range fcnmigrations.CodeMigrations {
		err := s.Add(ctx, m)
		if err != nil {
			panic(err)
		}
	}

	t, err := migrationsql.NewTarget(db)
	if err != nil {
		panic(err)
	}

	reporters.NewZapReporter(logger)

	_, err = migrations.Migrate(ctx, s, t, migrations.WithRunnerOptions(migrations.WithReporter(reporters.NewZapReporter(logger))))
	if err != nil {
		panic(err)
	}
}
