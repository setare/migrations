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

var _ = Describe("Resolver", func() {
	Describe("StepResolver", func() {
		It("should go forward without current migration", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			migration, err := migrations.StepResolver(source, target, 1).Resolve()
			Expect(err).ToNot(HaveOccurred())
			Expect(migration).To(Equal(m1))

			migration, err = migrations.StepResolver(source, target, 2).Resolve()
			Expect(err).ToNot(HaveOccurred())
			Expect(migration).To(Equal(m2))
		})

		It("should go forward without current migration", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))
			m3 := testingutils.NewMigration(time.Unix(2, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())
			Expect(source.Add(m3)).To(Succeed())

			target := testingutils.NewMemoryTarget()
			target.Add(m1)

			migration, err := migrations.StepResolver(source, target, 2).Resolve()
			Expect(err).ToNot(HaveOccurred())

			Expect(migration).To(Equal(m3))
		})

		It("should fail resolving a migration out of the upper boundary", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))
			m3 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())
			Expect(source.Add(m3)).To(Succeed())
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)

			// Try to resolve
			_, err := migrations.StepResolver(source, target, 2).Resolve()
			Expect(err).To(HaveOccurred())

			// Check error
			Expect(errors.Is(err, migrations.ErrStepOutOfIndex)).To(BeTrue())

		})

		It("should resolve a migration going backward", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))
			m3 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())
			Expect(source.Add(m3)).To(Succeed())
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)
			target.Add(m3)

			// Try to resolve
			migration, err := migrations.StepResolver(source, target, -2).Resolve()
			Expect(err).ToNot(HaveOccurred())
			Expect(migration).To(Equal(m1))
		})

		It("should fail resolving a migration out of the lower boundary", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))
			m3 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())
			Expect(source.Add(m3)).To(Succeed())
			target := testingutils.NewMemoryTarget()
			target.Add(m1)
			target.Add(m2)

			// Try to resolve
			_, err := migrations.StepResolver(source, target, -2).Resolve()
			Expect(err).To(HaveOccurred())

			// Check error
			Expect(errors.Is(err, migrations.ErrStepOutOfIndex)).To(BeTrue())

		})
	})

	Describe("MostRecentResolver", func() {
		It("should return the most recent migration", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			migration, err := migrations.MostRecentResolver(source).Resolve()
			Expect(err).ToNot(HaveOccurred())

			Expect(migration).To(Equal(m2))
		})

		It("should fail resolving when there is no migrations available", func() {
			// Prepare scenario
			source := migrations.NewSource()
			_, err := migrations.MostRecentResolver(source).Resolve()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})
	})

	Describe("MostRecentResolver", func() {
		It("should return the most recent migration", func() {
			// Prepare scenario
			source := migrations.NewSource()
			m1 := testingutils.NewMigration(time.Unix(0, 0))
			m2 := testingutils.NewMigration(time.Unix(1, 0))

			Expect(source.Add(m1)).To(Succeed())
			Expect(source.Add(m2)).To(Succeed())

			migration, err := migrations.FirstMigrationResolver(source).Resolve()
			Expect(err).ToNot(HaveOccurred())

			Expect(migration).To(Equal(m1))
		})

		It("should fail resolving when there is no migrations available", func() {
			// Prepare scenario
			source := migrations.NewSource()
			_, err := migrations.FirstMigrationResolver(source).Resolve()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNoMigrationsAvailable)).To(BeTrue())
		})
	})
})
