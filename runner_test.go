package migrations_test

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jamillosantos/migrations"
	"github.com/jamillosantos/migrations/testingutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	It("should execute a plan", func() {
		// # Prepare scenario
		// Create 5 normal migrations.
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(1, 0))
		m2 := testingutils.NewMigration(time.Unix(2, 0))
		m3 := testingutils.NewMigration(time.Unix(3, 0))
		m4 := testingutils.NewMigration(time.Unix(4, 0))
		m5 := testingutils.NewMigration(time.Unix(5, 0))

		target := testingutils.NewMemoryTarget()

		// Create an artificial plan simulating a migration
		plan := migrations.Plan{&migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m1,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m2,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m3,
		}, &migrations.Action{
			Action:    migrations.ActionTypeUndo,
			Migration: m4,
		}, &migrations.Action{
			Action:    migrations.ActionTypeUndo,
			Migration: m5,
		}}

		// Initialize and activate the runner
		runner := migrations.NewRunner(source, target)
		stats, err := runner.Execute(migrations.NewExecutionContext(nil), plan, nil)
		Expect(err).ToNot(HaveOccurred())

		// Check runner returned stats
		Expect(stats.Errored).To(BeEmpty())
		Expect(stats.Successful).To(HaveLen(5))
		Expect(stats.Successful[0].Migration).To(Equal(m1))
		Expect(stats.Successful[1].Migration).To(Equal(m2))
		Expect(stats.Successful[2].Migration).To(Equal(m3))
		Expect(stats.Successful[3].Migration).To(Equal(m4))
		Expect(stats.Successful[4].Migration).To(Equal(m5))

		// Check if the migrations were actually performed.
		migrationsDone, err := target.Done()
		Expect(err).ToNot(HaveOccurred())
		Expect(migrationsDone).To(ConsistOf(m1, m2, m3))

		Expect(m1.DoneCount).To(Equal(1))
		Expect(m2.DoneCount).To(Equal(1))
		Expect(m3.DoneCount).To(Equal(1))
		Expect(m4.UndoneCount).To(Equal(1))
		Expect(m5.UndoneCount).To(Equal(1))
	})

	It("should fail executing a migration", func() {
		// # Prepare scenario
		// Create 5 normal migrations.
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(1, 0))
		m2 := testingutils.NewMigration(time.Unix(2, 0))
		m3 := testingutils.NewMigration(time.Unix(3, 0))
		m4 := testingutils.NewMigration(time.Unix(4, 0))
		m5 := testingutils.NewMigration(time.Unix(5, 0))

		m3.DoErr = errors.New("forced error")

		target := testingutils.NewMemoryTarget()

		// Create an artificial plan simulating a migration
		plan := migrations.Plan{&migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m1,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m2,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m3,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m4,
		}, &migrations.Action{
			Action:    migrations.ActionTypeDo,
			Migration: m5,
		}}

		// Initialize and activate the runner
		runner := migrations.NewRunner(source, target)
		stats, err := runner.Execute(migrations.NewExecutionContext(nil), plan, nil)
		Expect(err).To(Equal(m3.DoErr))

		// Check runner returned stats
		Expect(stats.Errored).To(HaveLen(1))
		Expect(stats.Errored[0].Action).To(Equal(migrations.ActionTypeDo))
		Expect(stats.Errored[0].Migration).To(Equal(m3))
		Expect(stats.Successful).To(HaveLen(2))
		Expect(stats.Successful[0].Migration).To(Equal(m1))
		Expect(stats.Successful[1].Migration).To(Equal(m2))

		// Check if the migrations were actually performed.
		migrationsDone, err := target.Done()
		Expect(err).ToNot(HaveOccurred())
		Expect(migrationsDone).To(ConsistOf(m1, m2))

		Expect(m1.DoneCount).To(Equal(1))
		Expect(m2.DoneCount).To(Equal(1))
		Expect(m3.DoneCount).To(Equal(1))
		Expect(m4.DoneCount).To(Equal(0))
		Expect(m5.DoneCount).To(Equal(0))
	})

	It("should fail executing a plan when undoing a undoable migration", func() {
		// # Prepare scenario
		//
		// Create two migrations:
		// - An undoable migration
		// - A normal migration
		source := migrations.NewSource()
		m1 := testingutils.NewForwardMigration(time.Unix(1, 0)) // Cannot be undone
		m2 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m2)).To(Succeed())

		// Mark migration as executed by adding to the source.
		target := testingutils.NewMemoryTarget()
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())

		// Create an artificial plan undoing m1 that is "undoable".
		plan := migrations.Plan{
			&migrations.Action{
				Action:    migrations.ActionTypeUndo,
				Migration: m2,
			},
			&migrations.Action{
				Action:    migrations.ActionTypeUndo,
				Migration: m1,
			},
		}

		// Try to execute the plan.
		runner := migrations.NewRunner(source, target)
		_, err := runner.Execute(nil, plan, nil)

		// Get an `ErrMigrationNotUndoable` error
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, migrations.ErrMigrationNotUndoable)).To(BeTrue())
		vErr := err.(migrations.MigrationError)
		// Checks if the migration is expected one.
		Expect(vErr.Migration()).To(Equal(m1))

		// Ensures that no migration was performed and the system was not changed.
		migrationsDone, err := target.Done()
		Expect(err).ToNot(HaveOccurred())
		Expect(migrationsDone).To(ConsistOf(m1, m2))
	})
})
