package sql

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jamillosantos/migrations/v2"
	"github.com/jamillosantos/migrations/v2/sql/drivers"
)

var _ = Describe("SQL", func() {
	var (
		db     *sql.DB
		target *Target

		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		newDB, err := sql.Open("sqlite3", ":memory:")
		Expect(err).ToNot(HaveOccurred(), "should open the database")
		Expect(newDB.Ping()).To(Succeed(), "should ping the database")

		db = newDB

		newTarget, err := NewTarget(db, WithDriverOptions(drivers.WithDatabaseName("test")))
		Expect(err).ToNot(HaveOccurred(), "should create the target")

		target = newTarget

		Expect(target.Create(ctx)).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Close()).To(Succeed())
	})

	It("should add a new migration", func() {
		Expect(target.Add(ctx, "1")).ToNot(HaveOccurred())
		Expect(target.FinishMigration(ctx, "1")).ToNot(HaveOccurred())

		done, err := target.Done(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(done).To(Equal([]string{"1"}))
	})

	When("the migration does not exists", func() {
		It("should return an ErrMigrationNotFound", func() {
			Expect(target.Add(ctx, "1")).To(Succeed())
			Expect(target.Add(ctx, "1")).To(HaveOccurred()) // TODO(J): This should be handled by the driver. No SQLlite default support. (Only for testing)
		})
	})

	It("should remove the migration", func() {
		Expect(target.Add(ctx, "1")).ToNot(HaveOccurred())
		Expect(target.Remove(ctx, "1")).ToNot(HaveOccurred())

		done, err := target.Done(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(done).To(BeEmpty())
	})

	When("the migration does not exists", func() {
		It("should return an ErrMigrationNotFound", func() {
			Expect(target.Remove(ctx, "1")).To(MatchError(migrations.ErrMigrationNotFound))
		})
	})
})
