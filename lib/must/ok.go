package must

import (
	"fmt"
)

func OK[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(fmt.Errorf("no error assurance failed: %w", err))
	}

	return v
}
