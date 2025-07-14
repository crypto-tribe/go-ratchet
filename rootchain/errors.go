package rootchain

import (
	"errors"
)

var (
	// ErrApplyOptions is config options apply error.
	ErrApplyOptions = errors.New("apply options")

	// ErrCryptoAdvance is the chain crypto advance error.
	ErrCryptoAdvance = errors.New("advance")

	// ErrCryptoIsNil is an error when nil crypto was passed.
	ErrCryptoIsNil = errors.New("crypto is nil")

	// ErrKDF is the key derivation error.
	ErrKDF = errors.New("KDF")

	// ErrNewConfig is the config initialization error.
	ErrNewConfig = errors.New("new config")

	// ErrNewHasher is the hasher initialization error.
	ErrNewHasher = errors.New("new hasher")
)
