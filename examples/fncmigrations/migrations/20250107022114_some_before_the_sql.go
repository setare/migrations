package migrations

import (
	"context"
	"fmt"
)

var _ = Migration(func(ctx context.Context) error {
	fmt.Println("This happened before the SQL.")

	return nil
})
