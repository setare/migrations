package migrations

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type runnerObjs struct {
	runner *Runner
	source *MockSource
	target *MockTarget
	ctrl   *gomock.Controller
}

func createRunner(t *testing.T) (r runnerObjs) {
	t.Helper()
	r.ctrl = gomock.NewController(t)
	r.source = NewMockSource(r.ctrl)
	r.target = NewMockTarget(r.ctrl)
	r.runner = NewRunner(r.source, r.target)
	return
}

func TestNewRunner(t *testing.T) {
	ctrl := gomock.NewController(t)
	wantSource := NewMockSource(ctrl)
	wantTarget := NewMockTarget(ctrl)
	r := NewRunner(wantSource, wantTarget)
	assert.Equal(t, wantSource, r.source)
	assert.Equal(t, wantTarget, r.target)
}

func TestRunner_Execute(t *testing.T) {
	t.Run("should execute a plan", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "2")

		m1.EXPECT().Do(ctx).Return(nil)
		m2.EXPECT().CanUndo().Return(true)
		m2.EXPECT().Undo(ctx).Return(nil)

		s.target.EXPECT().Add(m1).Return(nil)
		s.target.EXPECT().Remove(m2).Return(nil)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeDo,
				Migration: m1,
			},
			&Action{
				Action:    ActionTypeUndo,
				Migration: m2,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.NoError(t, err)

		// Check runner returned stats
		assert.Empty(t, stats.Errored)
		assert.Len(t, stats.Successful, 2)
		assert.Equal(t, m1, stats.Successful[0].Migration)
		assert.Equal(t, m2, stats.Successful[1].Migration)
	})

	t.Run("should execute a plan with reporter", func(t *testing.T) {
		ctx := context.Background()
		s := createRunner(t)

		ctrl := gomock.NewController(t)

		reporter := NewMockRunnerReporter(s.ctrl)

		m1 := newMockMigration(ctrl, "1")
		m2 := newMockMigration(ctrl, "2")

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeDo,
				Migration: m1,
			},
			&Action{
				Action:    ActionTypeUndo,
				Migration: m2,
			},
		}

		m1.EXPECT().Do(ctx).Return(nil)
		m2.EXPECT().CanUndo().Return(true)
		m2.EXPECT().Undo(ctx).Return(nil)

		s.target.EXPECT().Add(m1).Return(nil)
		s.target.EXPECT().Remove(m2).Return(nil)

		reporter.EXPECT().BeforeExecute(ctx, plan).Return()
		reporter.EXPECT().AfterExecute(ctx, plan, gomock.Any(), nil).Return()

		reporter.EXPECT().BeforeExecuteMigration(ctx, ActionTypeDo, m1).Return()
		reporter.EXPECT().AfterExecuteMigration(ctx, ActionTypeDo, m1, nil).Return()

		reporter.EXPECT().BeforeExecuteMigration(ctx, ActionTypeUndo, m2).Return()
		reporter.EXPECT().AfterExecuteMigration(ctx, ActionTypeUndo, m2, nil).Return()

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, reporter)
		require.NoError(t, err)

		// Check runner returned stats
		assert.Empty(t, stats.Errored)
		assert.Len(t, stats.Successful, 2)
		assert.Equal(t, m1, stats.Successful[0].Migration)
		assert.Equal(t, m2, stats.Successful[1].Migration)
	})

	t.Run("should fail when trying to undo an undoable migration", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		m1.EXPECT().CanUndo().Return(false)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeUndo,
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		assert.ErrorIs(t, err, ErrMigrationNotUndoable)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		assert.Empty(t, stats.Errored)
	})

	t.Run("should fail when the migration Do fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		wantErr := errors.New("random error")

		m1.EXPECT().Do(gomock.Any()).Return(wantErr)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeDo,
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.ErrorIs(t, err, wantErr)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		require.Len(t, stats.Errored, 1)
	})

	t.Run("should fail when the target fails to Add the migration registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		wantErr := errors.New("random error")

		m1.EXPECT().Do(gomock.Any()).Return(nil)
		s.target.EXPECT().Add(gomock.Any()).Return(wantErr)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeDo,
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.ErrorIs(t, err, wantErr)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		require.Len(t, stats.Errored, 1)
	})

	t.Run("should fail when the migration Undo fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		wantErr := errors.New("random error")

		m1.EXPECT().CanUndo().Return(true)
		m1.EXPECT().Undo(gomock.Any()).Return(wantErr)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeUndo,
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.ErrorIs(t, err, wantErr)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		require.Len(t, stats.Errored, 1)
	})

	t.Run("should fail when the target fails to Remove the migration registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		wantErr := errors.New("random error")

		m1.EXPECT().Undo(gomock.Any()).Return(nil)
		m1.EXPECT().CanUndo().Return(true)
		s.target.EXPECT().Remove(gomock.Any()).Return(wantErr)

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionTypeUndo,
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.ErrorIs(t, err, wantErr)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		require.Len(t, stats.Errored, 1)
	})

	t.Run("should fail when the action is defined with an unknown type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		s := createRunner(t)

		m1 := newMockMigration(ctrl, "1")

		// Create an artificial plan simulating a migration
		plan := Plan{
			&Action{
				Action:    ActionType("unknown"),
				Migration: m1,
			},
		}

		// Initialize and activate the runner
		stats, err := s.runner.Execute(ctx, plan, nil)
		require.ErrorIs(t, err, ErrInvalidAction)
		require.NotNil(t, stats)
		assert.Empty(t, stats.Successful)
		require.Len(t, stats.Errored, 1)
	})
}
