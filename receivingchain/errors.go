package receivingchain

import (
	"errors"
	"fmt"

	cipher "golang.org/x/crypto/chacha20poly1305"
)

var (
	// ErrAddSkippedKey is the skipped key add error.
	ErrAddSkippedKey = errors.New("add skipped key")

	// ErrAdvanceChain is chain advance error.
	ErrAdvanceChain = errors.New("advance chain")

	// ErrApplyOptions is the config options apply error.
	ErrApplyOptions = errors.New("apply options")

	// ErrCryptoAdvanceChain is chain advance error from crypto provider.
	ErrCryptoAdvanceChain = errors.New("crypto advance chain")

	// ErrCryptoIsNil is the nil crypto error.
	ErrCryptoIsNil = errors.New("crypto is nil")

	// ErrDecodeHeader is the header decoding error.
	ErrDecodeHeader = errors.New("decode header")

	// ErrDecrypt is decryption error.
	ErrDecrypt = errors.New("decrypt")

	// ErrDecryptHeaderWithCurrentOrNextKey is the header decryption with current or next key error.
	ErrDecryptHeaderWithCurrentOrNextKey = errors.New("decrypt header with current or next key")

	// ErrDecryptHeaderWithCurrentKey is the header decryption with current key error.
	ErrDecryptHeaderWithCurrentKey = errors.New("decrypt header with current key")

	// ErrDecryptHeaderWithNextKey is the header decryption with next key error.
	ErrDecryptHeaderWithNextKey = errors.New("decrypt header with next key")

	// ErrDecryptMessage is the message decrypt error.
	ErrDecryptMessage = errors.New("decrypt message")

	// ErrDecryptWithSkippedKeys is decryption using skipped keys error.
	ErrDecryptWithSkippedKeys = errors.New("decrypt with skipped keys")

	// ErrDeleteSkippedKeys is the skipped keys deletion error.
	ErrDeleteSkippedKeys = errors.New("delete skipped keys")

	// ErrDeriveMessageCipherKeyAndNonce is the message cipher key and nonce derivation error.
	ErrDeriveMessageCipherKeyAndNonce = errors.New("derive message cipher key and nonce")

	// ErrGetSkippedKeysStorageIter is the skipped keys storage iterator obtaining error.
	ErrGetSkippedKeysStorageIter = errors.New("get skipped keys storage iter")

	// ErrHandleEncryptedHeader is the encrypted header handle error.
	ErrHandleEncryptedHeader = errors.New("handle encrypted header")

	// ErrHeaderKeyIsNil is the nil header key error.
	ErrHeaderKeyIsNil = errors.New("header key is nil")

	// ErrMasterKeyIsNil is the nil master key error.
	ErrMasterKeyIsNil = errors.New("master key is nil")

	// ErrNotEnoughEncryptedHeaderBytes is the not enough encrypted header bytes.
	ErrNotEnoughEncryptedHeaderBytes = fmt.Errorf(
		"encrypted header too shot, expected at least %d bytes",
		cipher.NonceSizeX+1,
	)

	// ErrNewCipher is the cipher initialization error.
	ErrNewCipher = errors.New("new cipher")

	// ErrNewConfig is the config initialization error.
	ErrNewConfig = errors.New("new config")

	// ErrNewHasher is the hasher initialization error.
	ErrNewHasher = errors.New("new hasher")

	// ErrOpenCipher is the cipher opening error.
	ErrOpenCipher = errors.New("open cipher")

	// ErrRatchet is the ratchet callback error.
	ErrRatchet = errors.New("ratchet")

	// ErrSkippedKeysNotFound is an error when skipped keys not found.
	ErrSkippedKeysNotFound = errors.New("skipped keys not found")

	// ErrSkippedKeysStorageIsNil is the nil skipped keys storage error.
	ErrSkippedKeysStorageIsNil = errors.New("skipped keys storage is nil")

	// ErrSkipCurrentChainKeys is the current chain keys skipping error.
	ErrSkipCurrentChainKeys = errors.New("skip current chain keys")

	// ErrSkipPreviousChainKeys is the previous chain keys skipping error.
	ErrSkipPreviousChainKeys = errors.New("skip previous chain keys")

	// ErrTooManySkippedMessageKeys is an error when there are too many skipped message keys.
	ErrTooManySkippedMessageKeys = errors.New("too many skipped message keys")

	// ErrWriteMasterKeyByteToMAC is the master key byte write error.
	ErrWriteMasterKeyByteToMAC = errors.New("write master key byte to MAC")

	// ErrWriteMessageKeyByteToMAC is the message key byte write error.
	ErrWriteMessageKeyByteToMAC = errors.New("write message key byte to MAC")
)
