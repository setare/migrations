package migrations

import (
	"context"
	"fmt"
	"os"
)

var _ = Migration(func(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "This will do some migration stuff.")

	return nil
})
