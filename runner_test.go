package migrations_test

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/setare/migrations"
	"github.com/setare/migrations/testingutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	It("should execute a plan", func() {
		// Prepare scenario
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(1, 0))
		m2 := testingutils.NewMigration(time.Unix(2, 0))
		m3 := testingutils.NewMigration(time.Unix(3, 0))
		m4 := testingutils.NewMigration(time.Unix(4, 0))
		m5 := testingutils.NewMigration(time.Unix(5, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m2)).To(Succeed())
		Expect(source.Add(m3)).To(Succeed())
		Expect(source.Add(m4)).To(Succeed())
		Expect(source.Add(m5)).To(Succeed())

		target := testingutils.NewMemoryTarget()
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())
		Expect(target.Add(m3)).To(Succeed())
		Expect(target.Add(m4)).To(Succeed())
		Expect(target.Add(m5)).To(Succeed())

		// Create the plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Reset()
		Expect(err).ToNot(HaveOccurred())

		runner := migrations.NewRunner(source, target)
		Expect(runner.Execute(nil, plan, nil)).To(Succeed())

		migrationsDone, err := target.Done()
		Expect(err).ToNot(HaveOccurred())

		Expect(migrationsDone).To(ConsistOf(m1, m2, m3, m4, m5))
	})

	It("should fail executing a plan when undoing a undoable migration", func() {
		// Prepare scenario
		source := migrations.NewSource()
		m1 := testingutils.NewForwardMigration(time.Unix(1, 0)) // Cannot be undone
		m2 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m2)).To(Succeed())

		target := testingutils.NewMemoryTarget()
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())

		// Create the plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Reset()
		Expect(err).ToNot(HaveOccurred())

		runner := migrations.NewRunner(source, target)
		_, err = runner.Execute(nil, plan, nil)
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, migrations.ErrMigrationNotUndoable)).To(BeTrue())
		vErr := err.(migrations.MigrationError)
		Expect(vErr.Migration()).To(Equal(m1))

		migrationsDone, err := target.Done()
		Expect(err).ToNot(HaveOccurred())

		Expect(migrationsDone).To(ConsistOf(m1))
	})
})
