package migrations

import (
	"context"
	"errors"
	"sort"
)

var ErrMigrationAlreadyExists = errors.New("migration already exists")

// Repository stores the migrations in memory and provide methods for acessing them.
type Repository struct {
	list []Migration
	byID map[string]Migration
}

func (r *Repository) List(_ context.Context) ([]Migration, error) {
	sort.Sort(&sortByID{r.list})
	return r.list, nil
}

type sortByID struct {
	list []Migration
}

func (s sortByID) Len() int {
	return len(s.list)
}

func (s sortByID) Less(i, j int) bool {
	return s.list[i].ID() < s.list[j].ID()
}

func (s sortByID) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (r *Repository) ByID(id string) (Migration, error) {
	if r.byID == nil {
		r.byID = make(map[string]Migration, 0)
		r.list = make([]Migration, 0)
	}
	if migration, ok := r.byID[id]; ok {
		return migration, nil
	}
	return nil, WrapMigrationID(ErrMigrationNotFound, id)
}

func (r *Repository) Add(migration Migration) error {
	if r.byID == nil {
		r.byID = make(map[string]Migration, 1)
		r.list = make([]Migration, 0, 1)
	}
	if _, ok := r.byID[migration.ID()]; ok {
		return WrapMigrationID(ErrMigrationAlreadyExists, migration.ID())
	}
	r.list = append(r.list, migration)
	r.byID[migration.ID()] = migration
	return nil
}
