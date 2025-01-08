package reporters

import (
	"context"

	"github.com/jamillosantos/migrations/v2"
)

type noOP struct {
}

func NewNoOP() migrations.RunnerReporter {
	return &noOP{}
}

func (n *noOP) BeforeExecute(_ context.Context, _ *migrations.BeforeExecuteInfo) {
}

func (n *noOP) BeforeExecuteMigration(_ context.Context, _ *migrations.BeforeExecuteMigrationInfo) {
}

func (n *noOP) AfterExecuteMigration(_ context.Context, _ *migrations.AfterExecuteMigrationInfo) {
}

func (n *noOP) AfterExecute(_ context.Context, _ *migrations.AfterExecuteInfo) {
}
