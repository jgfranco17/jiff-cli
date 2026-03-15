package errorhandling

import (
	"errors"
)

// ErrFailOnDiff is returned when differences are found and the fail flag is set.
var ErrFailOnDiff = errors.New("differences found and fail flag is set")
