package migrations

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createRepo() *Repository {
	return &Repository{
		list: make([]Migration, 0),
		byID: make(map[string]Migration),
	}
}

func newMockMigration(ctrl *gomock.Controller, id string) *MockMigration {
	m := NewMockMigration(ctrl)
	m.EXPECT().ID().Return(id).AnyTimes()
	m.EXPECT().String().Return(fmt.Sprintf("migration %s", id)).AnyTimes()
	return m
}

func Test_Repository_Add(t *testing.T) {
	t.Run("add a migration when source is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		wantMigration := newMockMigration(ctrl, "1")

		source := createRepo()

		err := source.Add(wantMigration)
		assert.NoError(t, err)

		require.Len(t, source.list, 1)
		require.Len(t, source.byID, 1)

		assert.Equal(t, wantMigration, source.list[0])
		assert.Contains(t, source.byID, wantMigration.ID())
	})

	t.Run("add a migration on the end of the list", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		m1 := newMockMigration(ctrl, "1")
		wantMigration := newMockMigration(ctrl, "2")

		source := createRepo()

		err := source.Add(m1)
		assert.NoError(t, err)

		err = source.Add(wantMigration)
		assert.NoError(t, err)

		require.Len(t, source.list, 2)
		require.Len(t, source.byID, 2)

		assert.Equal(t, m1, source.list[0])
		assert.Equal(t, wantMigration, source.list[1])
		assert.Contains(t, source.byID, m1.ID())
		assert.Contains(t, source.byID, wantMigration.ID())
	})

	t.Run("add a migration on the beginning of the list", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		m1 := newMockMigration(ctrl, "2")
		wantMigration := newMockMigration(ctrl, "1")

		source := createRepo()

		err := source.Add(m1)
		assert.NoError(t, err)

		err = source.Add(wantMigration)
		assert.NoError(t, err)

		require.Len(t, source.list, 2)
		require.Len(t, source.byID, 2)

		assert.Equal(t, m1, source.list[0])
		assert.Equal(t, wantMigration, source.list[1])
		assert.Contains(t, source.byID, m1.ID())
		assert.Contains(t, source.byID, wantMigration.ID())
	})

	t.Run("add a migration on the middle of the list", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		m1 := newMockMigration(ctrl, "1")
		wantMigration := newMockMigration(ctrl, "2")
		m3 := newMockMigration(ctrl, "3")

		source := createRepo()

		err := source.Add(m1)
		assert.NoError(t, err)
		err = source.Add(m3)
		assert.NoError(t, err)

		err = source.Add(wantMigration)
		assert.NoError(t, err)

		require.Len(t, source.list, 3)
		require.Len(t, source.byID, 3)

		assert.Equal(t, m1, source.list[0])
		assert.Equal(t, wantMigration, source.list[1])
		assert.Equal(t, m3, source.list[2])
		assert.Equal(t, wantMigration, source.list[1])
		assert.Contains(t, source.byID, m1.ID())
		assert.Contains(t, source.byID, wantMigration.ID())
		assert.Contains(t, source.byID, m3.ID())
	})

	t.Run("should fail adding a migration with same ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "1")

		source := createRepo()

		err := source.Add(m1)
		assert.NoError(t, err)

		err = source.Add(m2)
		assert.ErrorIs(t, err, ErrMigrationAlreadyExists)

		require.Len(t, source.list, 1)
		require.Len(t, source.byID, 1)
	})
}

func Test_Repository_List(t *testing.T) {
	t.Run("should list migrations in order", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := createRepo()
		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "2")
		m3 := newMockMigration(ctrl, "3")
		m4 := newMockMigration(ctrl, "4")
		repo.list = []Migration{m2, m1, m4, m3}
		gotMigrations, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, gotMigrations, 4)
		require.Equal(t, m1.ID(), gotMigrations[0].ID())
		require.Equal(t, m2.ID(), gotMigrations[1].ID())
		require.Equal(t, m3.ID(), gotMigrations[2].ID())
		require.Equal(t, m4.ID(), gotMigrations[3].ID())
	})

	t.Run("should fail when empty", func(t *testing.T) {
		ctx := context.Background()

		source := createRepo()
		gotMigrations, err := source.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, gotMigrations, 0)
	})
}

func Test_Repository_ByID(t *testing.T) {
	t.Run("should return the migration reference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		source := createRepo()
		m1 := newMockMigration(ctrl, "1")
		source.byID = map[string]Migration{
			"1": m1,
		}
		gotMigration, err := source.ByID("1")
		assert.NoError(t, err)
		assert.Equal(t, m1, gotMigration)
	})

	t.Run("should fail returning a non existing migration", func(t *testing.T) {
		source := createRepo()
		gotMigration, err := source.ByID("1")
		assert.ErrorIs(t, err, ErrMigrationNotFound)
		assert.Nil(t, gotMigration)
	})
}
