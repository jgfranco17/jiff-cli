package errorhandling

import (
	"errors"
)

var ErrFailOnDiff = errors.New("differences found and fail flag is set")
