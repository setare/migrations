# migrations

Migrations is an abstraction for a migration system and it can migrate anything.

## How to use

This example can be found in the `exapmles/postgres/main.go` file.

```go
package main

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jamillosantos/migrations/v2"
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
	defer func() {
		_ = logger.Sync()
    }()

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
```

## Generating new migration files

```bash
go run github.com/jamillosantos/migrations/v2/cli/migrations create -destination=migrations
```

This will create an emtpy migration file in the `migrations` with the correct timestamp and a description:

```
migrations/20210101000000_my_migration.sql`
```

## How it works

The `migrations` package is a simple abstraction for a migration system. It is able to migrate anything that migrations
would apply.

It has a simple design split into 3 main components:

### Source

A `Source` is the media that persists the migrations themselves. It is an entity that will load the migrations from SQL
files, or from a S3 bucket (or anything else you need!).

We have implemented support for the `fs.FS` interface, so you can use the `embed.FS` to load migrations from the binary.

Also, `fnc` was created to allow you to load migrations from a function AND they can be used together. Please, check the
`examples/fncmigrations` folder.

### Target

A `Target` is what the migrations persisted will be stored. If you are dealing with relational databases, like postgres,
you would use the `_migrations` table by default. However, you could create a JSON file that would store the executed 
migrations.

### Runner

The runner is the entity that will run the migrations. It will use the `Source` and `Target` to execute the migrations
to retrive information about the available migrations (from the `Source`) and what migrations are already applied (
from the `Target`).

## Extending

The `migrations` package is designed to be extended and you will, probably, only work with `Source`s and/or `Target`s.

### 1. Source

In some migration systems (like [golang-migrate/migrate](https://github.com/golang-migrate/migrate)) are stored as `.sql`
files and are stored into a directory as `<timestamp>_<description>.(up|down).sql`. In other systems, the migration is
a `func` that need to do some complex work and should connect many components before the actual database migration.

So, in the first example, the `Source` is a directory containing a bunch of .sql files with specific names. In the second
example, the `Source` are function that should be organized chronologically.

Hence, `Source` is the media that persists the migrations themselves. In practice, it is just an `interface{}` with a
bunch of methods that will list all available migrations ([check the code]()).

**TODO**: Link the Source interface.

### 2. Target

A `Target` are what the migrations are transforming. If you are dealing with relational databases, like postgres, you would
use a `TargetSQL` implementation (we provide one, check the our examples folder|**TODO**).

In practice, a `Target` is just an `interface{}` with a bunch of methods that will list executed migrations, mark and
unmark migrations as executed ([check the code]()).

**TODO**: Link the Target interface.

### 3. Executer

An Executer integrations `Source` and `Target` and is responsible for step actions, like `Do` and `Undo`. Each call will
step forward or backward one migration at a time.

### 4. Runner

Runners are, also, concrete. They capture the developer intentions and call the `Executer`.

Let's say that you want to migrate your system. By that, you mean to run all pending migrations. So the runner will
use the `Executer.Do` calling it multiple times to get all migrations executed.
