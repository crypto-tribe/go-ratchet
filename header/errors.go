package header

import (
	"errors"
)

// ErrNotEnoughBytes is an error when not enough bytes passed.
var ErrNotEnoughBytes = errors.New("not enough bytes")
