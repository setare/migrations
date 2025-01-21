package sql

import (
	"github.com/jamillosantos/migrations/v2/sql/drivers"
)

type targetOpts struct {
	driver        drivers.Driver
	tableName     string
	driverOptions []drivers.Option
}

type TargetOption func(target *targetOpts) error

func WithDriver(driver drivers.Driver) TargetOption {
	return func(target *targetOpts) error {
		target.driver = driver
		return nil
	}
}

func WithDriverOptions(options ...drivers.Option) TargetOption {
	return func(target *targetOpts) error {
		target.driverOptions = options
		return nil
	}
}

func WithTableName(tableName string) TargetOption {
	return func(target *targetOpts) error {
		target.tableName = tableName
		return nil
	}
}
