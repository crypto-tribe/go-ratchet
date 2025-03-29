package errlist

import "errors"

var (
	ErrCrypto             = errors.New("crypto")
	ErrInvalidValue       = errors.New("invalid value")
	ErrOption             = errors.New("option")
	ErrSkippedKeysStorage = errors.New("skipped keys storage")
)
