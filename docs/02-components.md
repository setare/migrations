---
sidebar_position: 2
slug: /components
---

# Components

## Source

Source is the component that will provide the migrations to be executed. It is
responsible for listing the migrations and loading the migration content.

For example, the [migration-sql](https://github.com/jamillosantos/migration-sql) package implements a source that can
load migrations from the filesystem or from a `go:embed`.

## Target

Target is the component responsible for listing the migrations that were executed and storing the migration execution.

## Planner

Planner is the component responsible for planning the migrations to be executed. It will receive the source and the
target and will plan the migrations to be executed.

## Runner

Runner is the component responsible for executing the migrations. It will receive the source and the target and a 
migration plan (generated from the Planner) and executes it.
