package migrations_test

import (
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/setare/migrations"
	"github.com/setare/migrations/code"
	"github.com/setare/migrations/testingutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planner", func() {
	It("should plan a from no current to the most recent", func() {
		// Prepare scenario
		source := code.NewSource()
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

		// Create the plan
		planner := migrations.NewPlanner(source, testingutils.NewMemoryTarget())
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.MostRecentResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
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
		// Prepare scenario
		source := code.NewSource()
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
		target.Add(m1)
		target.Add(m2)

		// Create plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.MostRecentResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
		Expect(plan).To(HaveLen(3))
		Expect(plan[0].Migration).To(Equal(m3))
		Expect(plan[0].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[1].Migration).To(Equal(m4))
		Expect(plan[1].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[2].Migration).To(Equal(m5))
		Expect(plan[2].Action).To(Equal(migrations.ActionTypeDo))
	})

	It("should plan a from the most recent migration to no migration", func() {
		// Prepare scenario
		source := code.NewSource()
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
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.FirstMigrationResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
		Expect(plan).To(HaveLen(5))
		Expect(plan[0].Migration).To(Equal(m5))
		Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[1].Migration).To(Equal(m4))
		Expect(plan[1].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[2].Migration).To(Equal(m3))
		Expect(plan[2].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[3].Migration).To(Equal(m2))
		Expect(plan[3].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[4].Migration).To(Equal(m1))
		Expect(plan[4].Action).To(Equal(migrations.ActionTypeUndo))
	})

	It("should plan a current migration to the most recent", func() {
		// Prepare scenario
		source := code.NewSource()
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
		target.Add(m1)
		target.Add(m2)
		target.Add(m3)

		// Create plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.FirstMigrationResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
		Expect(plan).To(HaveLen(3))
		Expect(plan[0].Migration).To(Equal(m3))
		Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[1].Migration).To(Equal(m2))
		Expect(plan[1].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[2].Migration).To(Equal(m1))
		Expect(plan[2].Action).To(Equal(migrations.ActionTypeUndo))
	})

	It("should undo all migrations then migrate all over again", func() {
		// Prepare scenario
		source := code.NewSource()
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
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.FirstMigrationResolver(source),
				migrations.MostRecentResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
		Expect(plan).To(HaveLen(10))
		Expect(plan[0].Migration).To(Equal(m5))
		Expect(plan[0].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[1].Migration).To(Equal(m4))
		Expect(plan[1].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[2].Migration).To(Equal(m3))
		Expect(plan[2].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[3].Migration).To(Equal(m2))
		Expect(plan[3].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[4].Migration).To(Equal(m1))
		Expect(plan[4].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[5].Migration).To(Equal(m1))
		Expect(plan[5].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[6].Migration).To(Equal(m2))
		Expect(plan[6].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[7].Migration).To(Equal(m3))
		Expect(plan[7].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[8].Migration).To(Equal(m4))
		Expect(plan[8].Action).To(Equal(migrations.ActionTypeDo))
		Expect(plan[9].Migration).To(Equal(m5))
		Expect(plan[9].Action).To(Equal(migrations.ActionTypeDo))
	})

	FIt("should migrate all trying to go to the most recent", func() {
		// Prepare scenario
		source := code.NewSource()
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

		// Create the plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.FirstMigrationResolver(source),
				migrations.MostRecentResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
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

	It("should do all migrations and then undo all", func() {
		// Prepare scenario
		source := code.NewSource()
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

		// Create the plan
		planner := migrations.NewPlanner(source, target)
		plan, err := planner.Plan(&migrations.PlanRequest{
			Resolvers: []migrations.MigrationResolver{
				migrations.MostRecentResolver(source),
				migrations.FirstMigrationResolver(source),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Check if the plan is correct
		//Expect(plan).To(HaveLen(10))
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
		Expect(plan[5].Migration).To(Equal(m5))
		Expect(plan[5].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[6].Migration).To(Equal(m4))
		Expect(plan[6].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[7].Migration).To(Equal(m3))
		Expect(plan[7].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[8].Migration).To(Equal(m2))
		Expect(plan[8].Action).To(Equal(migrations.ActionTypeUndo))
		Expect(plan[9].Migration).To(Equal(m1))
		Expect(plan[9].Action).To(Equal(migrations.ActionTypeUndo))
	})
})
