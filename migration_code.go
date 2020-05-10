package migrations

import "time"

type migrationCode struct {
	id          time.Time
	description string
	do          func(ExecutionContext) error
	undo        func(ExecutionContext) error
}

func (migration *migrationCode) ID() time.Time {
	return migration.id
}

func (migration *migrationCode) Description() string {
	return migration.description
}

func (migration *migrationCode) Do(executionContext ExecutionContext) error {
	return migration.do(executionContext)
}

func (migration *migrationCode) Undo(executionContext ExecutionContext) error {
	return migration.undo(executionContext)
}
