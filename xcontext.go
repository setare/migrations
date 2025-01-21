package migrations

import (
	"context"
	"time"
)

func detach(ctx context.Context) context.Context {
	return &detachedContext{ctx}
}

type detachedContext struct {
	context.Context
}

func (d *detachedContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (d *detachedContext) Done() <-chan struct{} {
	return nil
}
