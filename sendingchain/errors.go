package sendingchain

import (
	"errors"
)

var (
	// ErrAdvance is the chain advance error.
	ErrAdvance = errors.New("advance")

	// ErrApplyOptions is the config options apply error.
	ErrApplyOptions = errors.New("apply options")

	// ErrCryptoAdvance is an error to advance using crypto from config.
	ErrCryptoAdvance = errors.New("crypto advance")

	// ErrCryptoIsNil is an error when crypto is nil.
	ErrCryptoIsNil = errors.New("crypto is nil")

	// ErrDeriveMessageCipherKeyAndNonce is the key and nonce derivation error.
	ErrDeriveMessageCipherKeyAndNonce = errors.New("derive message cipher key and nonce")

	// ErrEncryptHeader is the header encryption error.
	ErrEncryptHeader = errors.New("encrypt header")

	// ErrEncryptMessage is the message encryption error.
	ErrEncryptMessage = errors.New("encrypt message")

	// ErrGenerateNonce is the nonce generation error.
	ErrGenerateNonce = errors.New("generate nonce")

	// ErrHeaderKeyIsNil is the header key nil error.
	ErrHeaderKeyIsNil = errors.New("header key is nil")

	// ErrMasterKeyIsNil is the master key nil error.
	ErrMasterKeyIsNil = errors.New("master key is nil")

	// ErrNewCipher is the cipher initialization error.
	ErrNewCipher = errors.New("new cipher")

	// ErrNewConfig is the config initialization error.
	ErrNewConfig = errors.New("new config")

	// ErrNewHasher is the hasher initialization error.
	ErrNewHasher = errors.New("new hasher")

	// ErrWriteMasterKeyByteToMAC is the master key byte write error.
	ErrWriteMasterKeyByteToMAC = errors.New("write master key byte to MAC")

	// ErrWriteMessageKeyByteToMAC is the message key byte write error.
	ErrWriteMessageKeyByteToMAC = errors.New("write message key byte to MAC")
)
