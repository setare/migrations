package migrations

import (
	"context"

	"github.com/jamillosantos/migrations/v2"
	"github.com/jamillosantos/migrations/v2/fnc"
)

var (
	CodeMigrations = make([]migrations.Migration, 0)
)

func Migration(do func(ctx context.Context) error) migrations.Migration {
	m := fnc.Migration(do, fnc.WithSkip(2))
	CodeMigrations = append(CodeMigrations, m)
	return m
}
