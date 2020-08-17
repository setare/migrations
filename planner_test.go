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

var _ = Describe("Planners", func() {
	Describe("MigratePLanner", func() {
		It("should fail planning from a empty Source", func() {
			// Create an empty source and target
			source := migrations.NewSource()
			target := testingutils.NewMemoryTarget()

			// Plan a migrate command.
			planner := migrations.MigratePlanner(source, target)
			_, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})

		It("should plan migration with no current migration", func() {
			// Create 5 normal migrations.
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

			// No current migration.

			// Plan a migrate command.
			planner := migrations.MigratePlanner(source, testingutils.NewMemoryTarget())
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1, m2, m3, m4 and m5.
			Expect(plan).To(HaveLen(5))
			Expect(plan[0].Migration).To(Equal(m1))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[1].Migration).To(Equal(m2))
			Expect(plan[1].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[2].Migration).To(Equal(m3))
			Expect(plan[2].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[3].Migration).To(Equal(m4))
			Expect(plan[3].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[4].Migration).To(Equal(m5))
			Expect(plan[4].Action).To(Equal(migrations.ActionTypeDo))
		})

		It("should plan a current migration to the most recent", func() {
			// Create 5 normal migrations.
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

			// Create a target and simulate migrating m1 and m2.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)

			// Plan a migrate command.
			planner := migrations.MigratePlanner(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m3, m4 and m5.
			Expect(plan).To(HaveLen(3))
			Expect(plan[0].Migration).To(Equal(m3))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[1].Migration).To(Equal(m4))
			Expect(plan[1].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[2].Migration).To(Equal(m5))
			Expect(plan[2].Action).To(Equal(migrations.ActionTypeDo))
		})

		It("should return an empty plan", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			// Simulate the migrate for all migrations
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)

			// Plan a migrate command.
			planner := migrations.MigratePlanner(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is empty: no migrations should be done.
			Expect(plan).To(BeEmpty())
		})

		It("should fail migrating when the current migration cannot be found in the source", func() {
			// Create 3 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			m3 := testingutils.NewMigration(time.Unix(3, 0))

			// Add only m1 and m2 to the source, leaving m3 out.
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			// Simulate the migrate m1 and m3.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m3)

			// Plan a migrate command.
			planner := migrations.MigratePlanner(source, target)
			plan, err := planner.Plan()

			// Planning expected to fail because the current migration was not found.
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrCurrentMigrationNotFound)).To(BeTrue())

			// Check if the plan is empty: no migrations should be done.
			Expect(plan).To(BeEmpty())
		})
	})

	Describe("RewindPlanner", func() {
		It("should fail planning from a empty Source", func() {
			// Create an empty source and target
			source := migrations.NewSource()
			target := testingutils.NewMemoryTarget()

			// Plan a rewind command.
			planner := migrations.RewindPlanner(source, target)
			_, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})

		It("should plan a rewind with no current migration (empty plan)", func() {
			// Create 5 normal migrations.
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

			// Plan a rewind command.
			planner := migrations.RewindPlanner(source, testingutils.NewMemoryTarget())
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1, m2, m3, m4 and m5.
			Expect(plan).To(BeEmpty())
		})

		It("should plan a rewind from the current migration", func() {
			// Create 5 normal migrations.
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

			// Create a target and simulate migrating m1 and m2.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)
			target.Add(m3)

			// Plan a rewind command.
			planner := migrations.RewindPlanner(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m3, m4 and m5.
			Expect(plan).To(HaveLen(3))
			Expect(plan[0].Migration).To(Equal(m3))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[1].Migration).To(Equal(m2))
			Expect(plan[1].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[2].Migration).To(Equal(m1))
			Expect(plan[2].Action).To(Equal(migrations.ActionTypeUndo))
		})

		It("should fail migrating when the current migration cannot be found in the source", func() {
			// Create 3 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			m3 := testingutils.NewMigration(time.Unix(3, 0))

			// Add only m1 and m2 to the source, leaving m3 out.
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			// Simulate the migrate m1 and m3.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m3)

			// Plan a rewind command.
			planner := migrations.RewindPlanner(source, target)
			plan, err := planner.Plan()

			// Planning expected to fail because the current migration was not found.
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrCurrentMigrationNotFound)).To(BeTrue())

			// Check if the plan is empty: no migrations should be done.
			Expect(plan).To(BeEmpty())
		})
	})

	Describe("ResetPlanner", func() {
		It("should fail planning from a empty Source", func() {
			// Create an empty source and target
			source := migrations.NewSource()
			target := testingutils.NewMemoryTarget()

			// Plan a rewind command.
			planner := migrations.ResetPlanner(source, target)
			_, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})

		It("should plan a reset with no current migration", func() {
			// Create 5 normal migrations.
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

			// No current migration

			// Plan a reset command.
			planner := migrations.ResetPlanner(source, testingutils.NewMemoryTarget())
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1, m2, m3, m4 and m5.
			Expect(plan).To(HaveLen(5))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[0].Migration).To(Equal(m1))
			Expect(plan[1].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[1].Migration).To(Equal(m2))
			Expect(plan[2].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[2].Migration).To(Equal(m3))
			Expect(plan[3].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[3].Migration).To(Equal(m4))
			Expect(plan[4].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[4].Migration).To(Equal(m5))
		})

		It("should plan a reset from the current migration", func() {
			// Create 5 normal migrations.
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

			// Create a target and simulate migrating m1, m2 and m3.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)
			target.Add(m3)

			// Plan a rewind command.
			planner := migrations.ResetPlanner(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: undo m1, m2 and m3 / do m1, m2, m3, m4 and m5.
			Expect(plan).To(HaveLen(8))
			Expect(plan[0].Migration).To(Equal(m3))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[1].Migration).To(Equal(m2))
			Expect(plan[1].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[2].Migration).To(Equal(m1))
			Expect(plan[2].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[3].Migration).To(Equal(m1))
			Expect(plan[3].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[4].Migration).To(Equal(m2))
			Expect(plan[4].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[5].Migration).To(Equal(m3))
			Expect(plan[5].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[6].Migration).To(Equal(m4))
			Expect(plan[6].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[7].Migration).To(Equal(m5))
			Expect(plan[7].Action).To(Equal(migrations.ActionTypeDo))
		})

		It("should fail migrating when the current migration cannot be found in the source", func() {
			// Create 3 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			m3 := testingutils.NewMigration(time.Unix(3, 0))

			// Add only m1 and m2 to the source, leaving m3 out.
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			// Simulate the migrate m1 and m3.
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m3)

			// Plan a rewind command.
			planner := migrations.ResetPlanner(source, target)
			plan, err := planner.Plan()

			// Planning expected to fail because the current migration was not found.
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrCurrentMigrationNotFound)).To(BeTrue())

			// Check if the plan is empty: no migrations should be done.
			Expect(plan).To(BeEmpty())
		})
	})

	Describe("StepPlanner", func() {
		It("should fail planning a step forward from a empty Source", func() {
			// Create an empty source and target
			source := migrations.NewSource()
			target := testingutils.NewMemoryTarget()

			// Plan a rewind command.
			planner := migrations.StepPlanner(1)(source, target)
			_, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})

		It("should fail planning a step backward from a empty Source", func() {
			// Create an empty source and target
			source := migrations.NewSource()
			target := testingutils.NewMemoryTarget()

			// Plan a rewind command.
			planner := migrations.StepPlanner(-1)(source, target)
			_, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})

		It("should plan a step forward with no current migration", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			// No current migration

			// Plan a reset command.
			planner := migrations.StepPlanner(1)(source, testingutils.NewMemoryTarget())
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1.
			Expect(plan).To(HaveLen(1))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[0].Migration).To(Equal(m1))
		})

		It("should plan a step forward with no current migration and a single migration in the source", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))

			Expect(source.Add(m1)).To(Succeed())

			// No current migration

			// Plan a reset command.
			planner := migrations.StepPlanner(1)(source, testingutils.NewMemoryTarget())
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1.
			Expect(plan).To(HaveLen(1))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
			Expect(plan[0].Migration).To(Equal(m1))
		})

		It("should plan a step backward with a single migration", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))

			Expect(source.Add(m1)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			target.Add(m1)

			// Plan a reset command.
			planner := migrations.StepPlanner(-1)(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: do m1.
			Expect(plan).To(HaveLen(1))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[0].Migration).To(Equal(m1))
		})

		It("should plan a step backward with no current migration", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)

			// Plan a reset command.
			planner := migrations.StepPlanner(-1)(source, target)
			plan, err := planner.Plan()
			Expect(err).ToNot(HaveOccurred())

			// Check if the plan is correct: undo m2.
			Expect(plan).To(HaveLen(1))
			Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
			Expect(plan[0].Migration).To(Equal(m2))
		})

		It("should plan an more step forward that available", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			target.Add(m1)

			// Try to plan 2 steps forward.
			planner := migrations.StepPlanner(2)(source, target)
			plan, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(migrations.ErrStepOutOfIndex))

			// Check if the plan is correct: <empty>.
			Expect(plan).To(BeEmpty())
		})

		It("should plan an more step backwards that available", func() {
			// Create 2 normal migrations.
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			target.Add(m1)

			// Try to plan 2 steps forward.
			planner := migrations.StepPlanner(-2)(source, target)
			plan, err := planner.Plan()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(migrations.ErrStepOutOfIndex))

			// Check if the plan is correct: <empty>.
			Expect(plan).To(BeEmpty())
		})
	})
})
