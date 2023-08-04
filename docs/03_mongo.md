---
sidebar_position: 3
---

# Mongo

| Repo:| [github.com/jamillosantos/migrations-mongo](https://github.com/jamillosantos/migrations-mongo) |
|------|------------------------------------------------------------------------------------------------|

The `migrations-mongo` implements a `Source` and `Target` for MongoDB.

## Creating a migration

First, you need to start creating your migrations. Below, you can see an example of a migration:

```go
// Filename: 20230803204512_add_index_to_users_updated_at.go

package migrations

import (
	"context"

	. "github.com/jamillosantos/migrations-fnc" // <- provides the Migration function.
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ = Migration(func(ctx context.Context) error {
	// DB should be a global variable initialized before the migration runs.
	c := DB.Collection("users") // Get the `users` collection.
	
	// Create the index.
	_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"updated_at": 1},
		Options: (&options.IndexOptions{}).
			SetName("idx_users_updated_at"),
	})
	return err
})
```

> As you can see, the `migrations-mongo` uses the `migrations-fnc` package to create the migrations.

You can have a file initialized with a different name pattern to declare the `DB` variable:

```go
// Filename: migrations/db.go

package migrations

import (
    "go.mongodb.org/mongo-driver/mongo"
)

var (
    // DB is the global database connection that will be used by the migrations.
    DB *mongo.Database
)
```

To run the migration you can:

```go
// ...

// Start a new mongo connection
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://guest:guest@localhost:27017"))
if err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "failed to connect with mongo: %s\n", err.Error())
    os.Exit(1)
}
// -----

db := client.Database("example")

// Initialize the DB that will be used in the migrations.
examplemigrations.DB = db

// Create the target where the migrations will run.
target, err := migrationsmongo.NewTarget(migrations.DefaultSource, db)
if err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "failed to initialize target: %s\n", err)
    os.Exit(1)
}

// runner is responsible for planning and running the migrations.
runner := migrations.NewRunner(migrations.DefaultSource, target)

reporter := migrationszap.NewRunnerReport(logger) // If you do not use zap, you can implement your own reporter.

// Run the migrations.
_, err = migrations.Migrate(ctx, runner, reporter)
if err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "failed to migrate: %s\n", err)
    os.Exit(1)
}

// ...
```


You can see the full example [here](https://github.com/jamillosantos/migrations-mongo/tree/main/internal/example).
