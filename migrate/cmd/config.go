package cmd

import "github.com/setare/migrations"

var (
	configFileFlag string
	yesFlag        bool
	planFlag       bool
)

var (
	planner *migrations.Planner
	runner  *migrations.Runner
)
