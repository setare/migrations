package code_test

import (
	"time"

	"github.com/setare/migrations/code"
	"github.com/setare/migrations/testingutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/setare/migrations"
)

var _ = Describe("Source", func() {
	Describe("Code", func() {
		It("should add a migration on an empty list", func() {
			s := code.NewSource()
			m := testingutils.NewMigration(time.Unix(0, 0))
			Expect(s.Add(m)).To(Succeed())
			l, err := s.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(l).To(HaveLen(1))
			Expect(m).To(Equal(l[0]))
			Expect(m.Next()).To(BeNil())
			Expect(m.Previous()).To(BeNil())
		})

		It("should sort migrations by ID while adding", func() {
			s := code.NewSource()
			m0 := testingutils.NewMigration(time.Unix(0, 0))
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(s.Add(m2)).To(Succeed())
			Expect(s.Add(m0)).To(Succeed())
			Expect(s.Add(m1)).To(Succeed())
			l, err := s.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(l).To(HaveLen(3))

			Expect(m0).To(Equal(l[0]))
			Expect(m0.Previous()).To(BeNil())
			Expect(m0.Next()).To(Equal(m1))

			Expect(m1).To(Equal(l[1]))
			Expect(m1.Previous()).To(Equal(m0))
			Expect(m1.Next()).To(Equal(m2))

			Expect(m2).To(Equal(l[2]))
			Expect(m2.Previous()).To(Equal(m1))
			Expect(m2.Next()).To(BeNil())
		})

		It("should fail adding an repeated ID", func() {
			s := code.NewSource()
			m := testingutils.NewMigration(time.Unix(0, 0))
			Expect(s.Add(m)).To(Succeed())
			err := s.Add(m)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, migrations.ErrNonUniqueMigrationID)).To(BeTrue())
		})

		It("should list available migrations", func() {
			s := code.NewSource()
			m0 := testingutils.NewMigration(time.Unix(0, 0))
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(s.Add(m0)).To(Succeed())
			Expect(s.Add(m1)).To(Succeed())
			Expect(s.Add(m2)).To(Succeed())
			l, err := s.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(l).To(HaveLen(3))
			Expect(m0).To(Equal(l[0]))
			Expect(m1).To(Equal(l[1]))
			Expect(m2).To(Equal(l[2]))
		})

		It("should get a migration by ID", func() {
			s := code.NewSource()
			m0 := testingutils.NewMigration(time.Unix(0, 0))
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(s.Add(m2)).To(Succeed())
			Expect(s.Add(m0)).To(Succeed())
			Expect(s.Add(m1)).To(Succeed())

			migration, err := s.ByID(m2.ID())
			Expect(err).ToNot(HaveOccurred())
			Expect(migration).To(Equal(m2))
		})

		It("should fail getting a migration by ID", func() {
			// Prepare scenario
			s := code.NewSource()
			m0 := testingutils.NewMigration(time.Unix(0, 0))
			m1 := testingutils.NewMigration(time.Unix(1, 0))
			m2 := testingutils.NewMigration(time.Unix(2, 0))
			Expect(s.Add(m2)).To(Succeed())
			Expect(s.Add(m0)).To(Succeed())
			Expect(s.Add(m1)).To(Succeed())

			// Tries to get an non existing migration by ID
			id := time.Unix(4, 0)
			migration, err := s.ByID(id)
			Expect(err).To(HaveOccurred())
			Expect(migration).To(BeNil())

			// Checks if the error returned is the expected
			Expect(errors.Is(err, migrations.ErrMigrationNotFound)).To(BeTrue())
			migrationIDErr, ok := err.(migrations.MigrationIDError)
			Expect(ok).To(BeTrue())
			Expect(migrationIDErr.MigrationID()).To(Equal(id))
		})
	})
})
