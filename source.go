package migrations

import (
	"context"
)

type memorySource struct {
	l []Migration
	m map[string]Migration
}

// NewMemorySource creates a source in which the migrations are stored only in memory. This is useful for
// `github.com/jamillosantos/migrations/v2/fnc` migrations.
func NewMemorySource() Source {
	return &memorySource{
		l: make([]Migration, 0),
		m: make(map[string]Migration),
	}
}

func (m *memorySource) Add(_ context.Context, migration Migration) error {
	if m.l == nil {
		m.l = make([]Migration, 0)
		m.m = make(map[string]Migration)
	}
	if _, ok := m.m[migration.ID()]; ok {
		return ErrMigrationAlreadyExists
	}
	m.l = append(m.l, migration)
	m.m[migration.ID()] = migration
	return nil
}

func (m *memorySource) Load(_ context.Context) (Repository, error) {
	return Repository{
		list: m.l,
		byID: m.m,
	}, nil
}
