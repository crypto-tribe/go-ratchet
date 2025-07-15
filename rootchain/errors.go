package rootchain

import (
	"errors"
)

var (
	// ErrAdvanceChain is the chain crypto advance error.
	ErrAdvanceChain = errors.New("advance")

	// ErrApplyOptions is config options apply error.
	ErrApplyOptions = errors.New("apply options")

	// ErrCryptoIsNil is an error when nil crypto was passed.
	ErrCryptoIsNil = errors.New("crypto is nil")

	// ErrKDF is the key derivation error.
	ErrKDF = errors.New("KDF")

	// ErrNewConfig is the config initialization error.
	ErrNewConfig = errors.New("new config")

	// ErrNewHasher is the hasher initialization error.
	ErrNewHasher = errors.New("new hasher")
)
