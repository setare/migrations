package migrations

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migratePlanner_Plan(t *testing.T) {
	t.Run("should migrate successfully", func(t *testing.T) {
		t.Run("should migrate from no migrations", func(t *testing.T) {
			// This case will test when there are no migrations run and the database is empty.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1")
			m2 := newMockMigration(ctrl, "2")
			m3 := newMockMigration(ctrl, "3")

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().Current().Return(nil, ErrNoCurrentMigration)

			source.EXPECT().List().Return([]Migration{
				m1, m2, m3,
			}, nil)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan()
			require.NoError(t, err)

			assert.Len(t, gotPlan, 3)
			assert.Equal(t, m1, gotPlan[0].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[0].Action)
			assert.Equal(t, m2, gotPlan[1].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[1].Action)
			assert.Equal(t, m3, gotPlan[2].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[2].Action)
		})

		t.Run("should migrate from already migrated migrations", func(t *testing.T) {
			// This case will test when there are a couple migrations already in place.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1") // previous migration
			m2 := newMockMigration(ctrl, "2") // current migration
			m3 := newMockMigration(ctrl, "3")
			m4 := newMockMigration(ctrl, "4") // target migration

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().Current().Return(m2, nil)

			source.EXPECT().List().Return([]Migration{
				m1, m2, m3, m4,
			}, nil)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan()
			require.NoError(t, err)

			assert.Len(t, gotPlan, 2)
			assert.Equal(t, m3, gotPlan[0].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[0].Action)
			assert.Equal(t, m4, gotPlan[1].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[1].Action)
		})

		t.Run("should return an empty plan when migrations is up to date", func(t *testing.T) {
			// This case will test when all migrations are already applied and there is nothing to be done.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1") // previous migration
			m2 := newMockMigration(ctrl, "2") // current migration
			m3 := newMockMigration(ctrl, "3")
			m4 := newMockMigration(ctrl, "4") // target migration

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().Current().Return(m4, nil)

			source.EXPECT().List().Return([]Migration{
				m1, m2, m3, m4,
			}, nil)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan()
			require.NoError(t, err)

			assert.Empty(t, gotPlan)
		})
	})

	t.Run("should fail when the listing migrations fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		wantErr := errors.New("random error")

		source.EXPECT().List().Return(nil, wantErr)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan()
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail we obtaining the current migration fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		wantErr := errors.New("random error")

		source.EXPECT().List().Return([]Migration{}, nil)

		target.EXPECT().Current().Return(nil, wantErr)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan()
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail when obtaining the current migration fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		m1 := newMockMigration(ctrl, "1")

		givenCurrentMigration := newMockMigration(ctrl, "2")

		source.EXPECT().List().Return([]Migration{
			m1,
		}, nil)

		target.EXPECT().Current().Return(givenCurrentMigration, nil)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan()
		assert.ErrorIs(t, err, ErrCurrentMigrationNotFound)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail when the current migration is more recent then the target migration", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		m1 := newMockMigration(ctrl, "1")

		givenCurrentMigration := newMockMigration(ctrl, "2")

		source.EXPECT().List().Return([]Migration{
			givenCurrentMigration,
			m1,
		}, nil)

		target.EXPECT().Current().Return(givenCurrentMigration, nil)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan()
		assert.ErrorIs(t, err, ErrCurrentMigrationMoreRecent)
		assert.Empty(t, gotPlan)
	})
}
