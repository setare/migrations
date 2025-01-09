package migrations

import (
	"context"
	"fmt"
	"os"
)

var _ = Migration(func(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "This happened before the SQL.")

	return nil
})
