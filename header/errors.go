package header

import (
	"errors"
)

var (
	// ErrNotEnoughBytes is an error when not enough bytes passed.
	ErrNotEnoughBytes = errors.New("not enough bytes")
)
