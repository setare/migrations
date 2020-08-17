package sql_test

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jamillosantos/migrations"
	migrationsSQL "github.com/jamillosantos/migrations/sql"
	"github.com/jamillosantos/migrations/testingutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL", func() {
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite3", ":memory:")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(db.Close()).To(Succeed())
	})

	It("should create the migrations table", func() {
		// Prepares scenario
		target, err := migrationsSQL.NewTarget(nil, db)
		Expect(err).ToNot(HaveOccurred())

		// Creates table
		Expect(target.Create()).To(Succeed())

		// Checks if the table exists
		var tableName string
		Expect(db.QueryRow("SELECT name FROM sqlite_master;").Scan(&tableName)).To(Succeed())
		Expect(tableName).To(Equal("_migrations"))
	})

	It("should create the migrations table customizing the table name", func() {
		// Prepares scenario
		target, err := migrationsSQL.NewTarget(nil, db, migrationsSQL.Table("new_migration_table"))
		Expect(err).ToNot(HaveOccurred())

		// Creates table
		Expect(target.Create()).To(Succeed())

		// Checks if the table exists
		var tableName string
		Expect(db.QueryRow("SELECT name FROM sqlite_master;").Scan(&tableName)).To(Succeed())
		Expect(tableName).To(Equal("new_migration_table"))
	})

	It("should receive the error from an option", func() {
		// Prepares scenario
		errOpt := errors.New("error from option")
		target, err := migrationsSQL.NewTarget(nil, db, migrationsSQL.OptError(errOpt))
		Expect(err).To(HaveOccurred())
		Expect(target).To(BeNil())
		Expect(err).To(Equal(errOpt))
	})

	It("should destroy the migrations table", func() {
		// Prepares scenario
		target, err := migrationsSQL.NewTarget(nil, db)
		Expect(err).ToNot(HaveOccurred())
		Expect(target.Create()).To(Succeed())

		// Destroys table
		Expect(target.Destroy()).To(Succeed())

		// Checks if the table exists.
		var tableName string
		err = db.QueryRow("SELECT name FROM sqlite_master;").Scan(&tableName)
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(sql.ErrNoRows))
	})

	It("should destroy the migrations table with a customized name", func() {
		// Prepares scenario
		target, err := migrationsSQL.NewTarget(nil, db, migrationsSQL.Table("new_migration_table"))
		Expect(err).ToNot(HaveOccurred())
		Expect(target.Create()).To(Succeed())

		// Destroys table
		Expect(target.Destroy()).To(Succeed())

		// Checks if the table exists.
		var tableName string
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table';").Scan(&tableName)
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(sql.ErrNoRows))
	})

	It("should add migrations as executed", func() {
		// Prepares scenarios
		target, err := migrationsSQL.NewTarget(nil, db)
		Expect(err).ToNot(HaveOccurred())
		Expect(target.Create()).To(Succeed())

		// Creates and adds migrations
		m1 := testingutils.NewMigration(time.Unix(0, 0))
		m2 := testingutils.NewMigration(time.Unix(1, 0))
		m3 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(target.Add(m3)).To(Succeed())
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())

		// Check the database for the tables
		rs, err := db.Query("SELECT id FROM _migrations ORDER BY id;")
		Expect(err).ToNot(HaveOccurred())
		defer rs.Close()

		var id int64

		// finds m1
		Expect(rs.Next()).To(BeTrue())
		Expect(rs.Scan(&id)).To(Succeed())
		Expect(id).To(Equal(m1.ID().Unix()))

		// finds m2
		Expect(rs.Next()).To(BeTrue())
		Expect(rs.Scan(&id)).To(Succeed())
		Expect(id).To(Equal(m2.ID().Unix()))

		// finds m3
		Expect(rs.Next()).To(BeTrue())
		Expect(rs.Scan(&id)).To(Succeed())
		Expect(id).To(Equal(m3.ID().Unix()))

		// EOF
		Expect(rs.Next()).To(BeFalse())
	})

	It("should remove migrations", func() {
		// Prepare scenario
		m1 := testingutils.NewMigration(time.Unix(0, 0))
		m2 := testingutils.NewMigration(time.Unix(1, 0))
		m3 := testingutils.NewMigration(time.Unix(2, 0))

		target, err := migrationsSQL.NewTarget(nil, db)
		Expect(err).ToNot(HaveOccurred())

		Expect(target.Create()).To(Succeed())
		Expect(target.Add(m3)).To(Succeed())
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())

		// Removes the migration m3
		Expect(target.Remove(m3)).To(Succeed())

		// Check the database for the tables
		rs, err := db.Query("SELECT id FROM _migrations ORDER BY id;")
		Expect(err).ToNot(HaveOccurred())
		defer rs.Close()

		var id int64

		// Find m1
		Expect(rs.Next()).To(BeTrue())
		Expect(rs.Scan(&id)).To(Succeed())
		Expect(id).To(Equal(m1.ID().Unix()))

		// Find m2
		Expect(rs.Next()).To(BeTrue())
		Expect(rs.Scan(&id)).To(Succeed())
		Expect(id).To(Equal(m2.ID().Unix()))

		// EOF
		Expect(rs.Next()).To(BeFalse())
	})

	It("should get the current migration", func() {
		// Prepare scenario
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(0, 0))
		m2 := testingutils.NewMigration(time.Unix(1, 0))
		m3 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m2)).To(Succeed())
		Expect(source.Add(m3)).To(Succeed())

		target, err := migrationsSQL.NewTarget(source, db)
		Expect(err).ToNot(HaveOccurred())

		Expect(target.Create()).To(Succeed())
		Expect(target.Add(m1)).To(Succeed())

		// Get the current migration
		currentMigration, err := target.Current()
		Expect(err).ToNot(HaveOccurred())
		Expect(currentMigration).To(Equal(m1))

		// Update the current migration
		Expect(target.Add(m3)).To(Succeed())

		// Get the current migration²
		currentMigration, err = target.Current()
		Expect(err).ToNot(HaveOccurred())
		Expect(currentMigration).To(Equal(m3))

		// Add a migration that will not change the current migration
		Expect(target.Add(m2)).To(Succeed())

		// Get the current migration³
		currentMigration, err = target.Current()
		Expect(err).ToNot(HaveOccurred())
		Expect(currentMigration).To(Equal(m3))
	})

	It("should list of the executed migrations", func() {
		// Prepare scenario
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(0, 0))
		m2 := testingutils.NewMigration(time.Unix(1, 0))
		m3 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m2)).To(Succeed())
		Expect(source.Add(m3)).To(Succeed())

		target, err := migrationsSQL.NewTarget(source, db)
		Expect(err).ToNot(HaveOccurred())

		Expect(target.Create()).To(Succeed())
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())
		Expect(target.Add(m3)).To(Succeed())

		// List executed migrations
		list, err := target.Done()
		Expect(err).ToNot(HaveOccurred())
		Expect(list).To(HaveLen(3))

		Expect(list[0]).To(Equal(m1))
		Expect(list[1]).To(Equal(m2))
		Expect(list[2]).To(Equal(m3))
	})

	It("should fail listing a migration that is not in the source", func() {
		// Prepare scenario
		source := migrations.NewSource()
		m1 := testingutils.NewMigration(time.Unix(0, 0))
		m2 := testingutils.NewMigration(time.Unix(1, 0))
		m3 := testingutils.NewMigration(time.Unix(2, 0))

		Expect(source.Add(m1)).To(Succeed())
		Expect(source.Add(m3)).To(Succeed())

		target, err := migrationsSQL.NewTarget(source, db)
		Expect(err).ToNot(HaveOccurred())

		Expect(target.Create()).To(Succeed())
		Expect(target.Add(m1)).To(Succeed())
		Expect(target.Add(m2)).To(Succeed())
		Expect(target.Add(m3)).To(Succeed())

		// List executed migrations
		_, err = target.Done()
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, migrations.ErrMigrationNotFound)).To(BeTrue())
	})
})
