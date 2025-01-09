package drivers

import (
	"database/sql/driver"
	"errors"
	"reflect"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
)

var drivers = map[string]DriverConstructor{
	"sql":        newSQL,
	"postgres":   newPostgres,
	"*pq.Driver": newPostgres,
}

// Register will register a new driver constructor for the given driver name.
func Register(name string, constructor DriverConstructor) {
	drivers[name] = constructor
}

type DBWithDriver interface {
	DB
	Driver() driver.Driver
}

func DriverFromDB(db DBWithDriver, options ...Option) (Driver, error) {
	driverName := reflect.TypeOf(db.Driver()).String()
	constructor, ok := drivers[driverName]
	if !ok {
		return newSQL(db, options...)
	}
	return constructor(db, options...)
}
