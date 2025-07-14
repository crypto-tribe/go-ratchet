package chainscommon

import (
	"errors"
)

var (
	// ErrKDF is an error to derive keys.
	ErrKDF = errors.New("KDF")

	// ErrNewHasher is an error to create a new hasher.
	ErrNewHasher = errors.New("new hasher")
)
