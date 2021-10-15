package migrations

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSource(t *testing.T) {
	source := NewSource()
	assert.IsType(t, &ListSource{}, source)

	assert.NotNil(t, source.byID)
	assert.NotNil(t, source.list)
}

func createBaseSource() *ListSource {
	return &ListSource{
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

func Test_ListSource_Add(t *testing.T) {
	t.Run("add a migration when source is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		wantMigration := newMockMigration(ctrl, "1")

		wantMigration.EXPECT().SetPrevious(nil)
		wantMigration.EXPECT().SetNext(nil)

		source := createBaseSource()

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

		m1.EXPECT().SetPrevious(nil)
		m1.EXPECT().SetNext(nil)

		wantMigration.EXPECT().SetPrevious(m1)
		wantMigration.EXPECT().SetNext(nil)

		m1.EXPECT().SetNext(wantMigration)
		wantMigration.EXPECT().Previous().Return(m1)

		source := createBaseSource()

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

		m1.EXPECT().SetPrevious(nil)
		m1.EXPECT().SetNext(nil)

		m1.EXPECT().Previous().Return(nil)
		wantMigration.EXPECT().SetPrevious(nil)
		wantMigration.EXPECT().Previous().Return(nil)
		wantMigration.EXPECT().SetNext(m1)
		m1.EXPECT().SetPrevious(wantMigration)

		source := createBaseSource()

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

		gomock.InOrder(
			// Adding m1
			m1.EXPECT().SetPrevious(nil),
			m1.EXPECT().SetNext(nil),

			// Adding m3
			m3.EXPECT().SetPrevious(m1),
			m3.EXPECT().Previous().Return(m1),
			m1.EXPECT().SetNext(m3),
			m3.EXPECT().SetNext(nil),

			// Adding wantMigration in between m1 and m3
			m3.EXPECT().Previous().Return(m1),
			wantMigration.EXPECT().SetPrevious(m1),
			wantMigration.EXPECT().Previous().Return(m1),
			m1.EXPECT().SetNext(wantMigration),
			wantMigration.EXPECT().SetNext(m3),
			m3.EXPECT().SetPrevious(wantMigration),
		)

		// ----------

		source := createBaseSource()

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

		m1.EXPECT().SetPrevious(nil)
		m1.EXPECT().SetNext(nil)

		source := createBaseSource()

		err := source.Add(m1)
		assert.NoError(t, err)

		err = source.Add(m2)
		assert.ErrorIs(t, err, ErrNonUniqueMigrationID)

		require.Len(t, source.list, 1)
		require.Len(t, source.byID, 1)
	})
}

func Test_ListSource_List(t *testing.T) {
	t.Run("should list migrations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		source := createBaseSource()
		m1 := newMockMigration(ctrl, "1")
		source.list = []Migration{m1}
		gotMigrations, err := source.List()
		assert.NoError(t, err)
		assert.Len(t, gotMigrations, 1)
	})

	t.Run("should fail when empty", func(t *testing.T) {
		source := createBaseSource()
		gotMigrations, err := source.List()
		assert.ErrorIs(t, err, ErrNoMigrationsAvailable)
		assert.Len(t, gotMigrations, 0)
	})
}

func Test_ListSource_ByID(t *testing.T) {
	t.Run("should return the migration reference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		source := createBaseSource()
		m1 := newMockMigration(ctrl, "1")
		source.byID = map[string]Migration{
			"1": m1,
		}
		gotMigration, err := source.ByID("1")
		assert.NoError(t, err)
		assert.Equal(t, m1, gotMigration)
	})

	t.Run("should fail returning a non existing migration", func(t *testing.T) {
		source := createBaseSource()
		gotMigration, err := source.ByID("1")
		assert.ErrorIs(t, err, ErrMigrationNotFound)
		assert.Nil(t, gotMigration)
	})
}
