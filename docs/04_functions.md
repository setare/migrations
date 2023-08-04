---
sidebar_position: 4
---

# Functions

`| Repo:| [github.com/jamillosantos/migrations-fnc](https://github.com/jamillosantos/migrations-fnc) |
|------|---------------------------------------------|`

The `migrations-fnc` implements a `Migration` that receives a function. Whenever the migration is executed, the function
is called.

With this, you can basically migrate anything. This can be combined with any `Tartget` implementation. So, you can 
migrate whatever you want and store the migration state in a SQL or Mongo database, for example.

## How a "function-migration" will look like?

The implementation of a migration is very straight forward:

```go
// 20231015015442_create_customers_table.go

package migrations

import (
	"context"

	. "github.com/jamillosantos/migrations-fnc"
)

var _ = Migration(func(ctx context.Context) error {
	fmt.Println("Execute a SQL or a MongoDB index creation!")
	return nil
})
```

The `Migration` method will use the `runtime.Caller` and discover the name of the file. This name will be used to 
extract both `ID` and `description`. For the given example: `20231015015442` will be the ID and `create customers table`
will be the description.

## What do migrations look like in my project?

For the function migrations, you can write migrations as `.go` files that will be compiled into the binary. So, in your
project that can be a `migrations` directory:

```
./migrations
├── 20211015015442_create_customers_table.go
├── 20211015015556_add_age_to_customers_table.go
└── 20211015045556_add_birthday_to_customers_table.go
```

You need to make sure the `migrations` package is, somehow, imported in your `main.go` file. If you are NOT using the
`migrations` package, you can import it as `_ "<project>/migrations"`:

```go
package main

import (
	// ...
	_ "<project>/migrations"
	// ...
)

func main () {
	// ...
}
```

> The [migrations-mongo](/mongo) uses the `migration-fnc` to create Mongo migrations. 