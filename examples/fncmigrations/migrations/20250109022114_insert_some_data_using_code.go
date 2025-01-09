package migrations

import (
	"context"
	"fmt"
)

var _ = Migration(func(ctx context.Context) error {
	fmt.Println("This will do some migration stuff.")

	return nil
})
