package migrations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_migratePlanner_Plan(t *testing.T) {
	t.Run("should migrate successfully", func(t *testing.T) {
		t.Run("should migrate from no migrations", func(t *testing.T) {
			ctx := context.Background()

			// This case will test when there are no migrations run and the database is empty.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1")
			m2 := newMockMigration(ctrl, "2")
			m3 := newMockMigration(ctrl, "3")

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().
				Current(ctx).
				Return("", ErrNoCurrentMigration)

			source.EXPECT().
				Load(ctx).
				Return(
					RepositoryBuilder().
						WithMigration(
							m1, m2, m3,
						).Build(),
					nil,
				)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan(ctx)
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
			ctx := context.Background()

			// This case will test when there are a couple migrations already in place.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1") // previous migration
			m2 := newMockMigration(ctrl, "2") // current migration
			m3 := newMockMigration(ctrl, "3")
			m4 := newMockMigration(ctrl, "4") // target migration

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().
				Current(ctx).
				Return(m2.ID(), nil)

			target.EXPECT().
				Done(ctx).
				Return([]string{
					m1.ID(),
				}, nil)

			source.EXPECT().
				Load(ctx).
				Return(
					RepositoryBuilder().WithMigration(
						m1, m2, m3, m4,
					).Build(),
					nil,
				)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan(ctx)
			require.NoError(t, err)

			assert.Len(t, gotPlan, 2)
			assert.Equal(t, m3, gotPlan[0].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[0].Action)
			assert.Equal(t, m4, gotPlan[1].Migration)
			assert.Equal(t, ActionTypeDo, gotPlan[1].Action)
		})

		t.Run("should return an empty plan when migrations is up to date", func(t *testing.T) {
			ctx := context.Background()

			// This case will test when all migrations are already applied and there is nothing to be done.
			ctrl := gomock.NewController(t)

			m1 := newMockMigration(ctrl, "1") // previous migration
			m2 := newMockMigration(ctrl, "2") // current migration
			m3 := newMockMigration(ctrl, "3")
			m4 := newMockMigration(ctrl, "4") // target migration

			source := NewMockSource(ctrl)
			target := NewMockTarget(ctrl)

			target.EXPECT().
				Current(ctx).
				Return(m4.ID(), nil)

			target.EXPECT().
				Done(ctx).
				Return([]string{
					m1.ID(), m2.ID(), m3.ID(), m4.ID(),
				}, nil)

			source.EXPECT().
				Load(ctx).
				Return(
					RepositoryBuilder().
						WithMigration(m1, m2, m3, m4).
						Build(),
					nil,
				)

			planner := MigratePlanner(source, target)
			gotPlan, err := planner.Plan(ctx)
			require.NoError(t, err)

			assert.Empty(t, gotPlan)
		})
	})

	t.Run("should fail when the listing migrations fail", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		wantErr := errors.New("random error")

		source.EXPECT().
			Load(ctx).
			Return(Repository{}, wantErr)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan(ctx)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail we obtaining the current migration fails", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		wantErr := errors.New("random error")

		source.EXPECT().
			Load(ctx).
			Return(Repository{}, nil)

		target.EXPECT().
			Current(ctx).
			Return("", wantErr)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan(ctx)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail we obtaining list of migrations applied fails", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		wantErr := errors.New("random error")

		m1 := newMockMigration(ctrl, "1")

		source.EXPECT().
			Load(ctx).
			Return(
				RepositoryBuilder().
					WithMigration(
						m1,
					).Build(),
				nil,
			)

		target.EXPECT().
			Current(ctx).
			Return(m1.ID(), nil)

		target.EXPECT().
			Done(ctx).
			Return([]string{}, wantErr)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan(ctx)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail when a migration is applied but not listed", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "2")
		m3 := newMockMigration(ctrl, "3")

		source.EXPECT().
			Load(ctx).
			Return(
				RepositoryBuilder().
					WithMigration(
						m1,
						m3,
					).Build(),
				nil,
			)

		target.EXPECT().
			Current(ctx).
			Return(m3.ID(), nil)

		target.EXPECT().
			Done(ctx).
			Return([]string{
				m1.ID(),
				m2.ID(),
				m3.ID(),
			}, nil)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan(ctx)
		assert.ErrorIs(t, err, ErrMigrationNotFound)
		nErr, ok := err.(MigrationIDError)
		require.True(t, ok)
		assert.Equal(t, m2.ID(), nErr.MigrationID())
		assert.Empty(t, gotPlan)
	})

	t.Run("should fail when a stale migration is detected", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)

		source := NewMockSource(ctrl)
		target := NewMockTarget(ctrl)

		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "2")
		m3 := newMockMigration(ctrl, "3")

		source.EXPECT().
			Load(ctx).
			Return(
				RepositoryBuilder().
					WithMigration(
						m1,
						m2,
						m3,
					).Build(),
				nil,
			)

		target.EXPECT().
			Current(ctx).
			Return(m3.ID(), nil)

		target.EXPECT().
			Done(ctx).
			Return([]string{
				m1.ID(),
				m3.ID(),
			}, nil)

		planner := MigratePlanner(source, target)
		gotPlan, err := planner.Plan(ctx)
		assert.ErrorIs(t, err, ErrStaleMigrationDetected)
		assert.Empty(t, gotPlan)
	})
}
