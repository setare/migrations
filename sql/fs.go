package sql

import (
	"io/fs"
	"os"

	"github.com/jamillosantos/migrations/v2"
)

// SourceFromFS creates a new source based on the provided fs.ReadDirFS and folder.
// This is useful for creating a source based on a virtual filesystem (for example `go:embed`).
// Example:
//
//	//go:embed migrations
//	var migrationsFS embed.FS
//
//	source, err := sql.SourceFromFS(dbGetter, migrationsFS, "migrations")
func SourceFromFS(dbGetter func() DBExecer, fs fs.ReadDirFS, folder string) (migrations.Source, error) {
	return &source{
		dbGetter: dbGetter,

		fs:     fs,
		folder: folder,
	}, nil
}

// SourceFromDirectory creates a new source based on the provided folder in the disk.
func SourceFromDirectory(dbGetter func() DBExecer, folder string) (migrations.Source, error) {
	return &source{
		dbGetter: dbGetter,

		fs:     os.DirFS(folder).(fs.ReadDirFS),
		folder: folder,
	}, nil
}
