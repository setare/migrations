---
sidebar_position: 5
---

# Creating Migrations

For creating migrations you can use the CLI provided by the `migrations-sql` package.

## Usage

### Using Without Installation

If you want to keep your `migrations-sql` CLI sync with your project dependencies
you can take advantage of the Go package infrastructure and add the following file:


```go title="tools/tools.go"
package tools

import (
    _ "github.com/jamillosantos/migrations-sql/cli/migrations-sql"
)
```

After that, you just need to use the command as follow:

```shell
go run github.com/jamillosantos/migrations-sql/cli/migrations-sql create
```

This approach has some advantages:
1. First, you don't need to install the `migrations-sql`, so no new binaries on
   your system.

2. Second, if you use a tool like dependabot or renovate to keep your dependencies
    up to date, `migrations-sql` will be upgrade automatically when you go.mod is
    updated.

### Installing the CLI

If you want to have the `migrations-sql` CLI installed on your system, you can do:

```shell
$ go install github.com/jamillosantos/migrations-sql/cli/migrations-sql@latest
```

Now, you should be able to run `migrations-sql` command from anywhere on you computer.

# CLI Options

## Create

```shell
migrations-sql create [--destination=<destination>] [--extension=sql] [description] [flags]
```

| Flag                                       | Description                                                                         |
|:-------------------------------------------|:------------------------------------------------------------------------------------|
| `description`                              | The description of the migration. This will be used to generate the migration name. If no description is given, the command will ask for one interactively. |
| `--destination=<folder>`<br/>`-d <folder>` | The destination folder where the new migration will be created. Defaults to `.`.    |
| `--extension=sql`<br/>`-e sql`             | The extension of the migration file. Defaults to `sql`.                             |
| `--undo`                                   | Flag that enables the generation of the undo file. Defaults to `false`.             |
| `--down`                                   | Flag that enables the generation of the down file (same as undo, but different naming). Defaults to `false`. |
| `--format=<format>`                        | Format of the migration ID using the Go time.Date.Format standard. Defaults to `20060102150405`. (`unix` alias for unix timestamp formatting) |

### Example 1:

```shell
go run github.com/jamillosantos/migrations-sql/cli/migrations-sql -d migrations create customer table
```

The command above will create the `migrations/create_customer_table.sql` file.

### Example 2:

```shell
go run github.com/jamillosantos/migrations-sql/cli/migrations-sql -d migrations --undo create customer table
```

The command above will create the `migrations/create_customer_table.do.sql` and `migrations/create_customer_table.undo.sql` files.

Similarly, you can use the `--down` flag to generate `.down` instead of `.undo` files.