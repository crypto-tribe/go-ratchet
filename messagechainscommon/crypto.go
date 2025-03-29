package messagechainscommon

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/platform-inf/go-ratchet/keys"
)

const cryptoMessageCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX

var (
	cryptoMessageCipherKDFSalt = make([]byte, cryptoMessageCipherKDFOutputLen)
	cryptoMessageCipherKDFInfo = []byte("message cipher")
)

func DeriveMessageCipherKeyAndNonce(messageKey keys.Message) ([]byte, []byte, error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash
		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	kdf := hkdf.New(getHasher, messageKey.Bytes, cryptoMessageCipherKDFSalt, cryptoMessageCipherKDFInfo)

	output := make([]byte, cryptoMessageCipherKDFOutputLen)
	if _, err := io.ReadFull(kdf, output); err != nil {
		return nil, nil, fmt.Errorf("KDF: %w", err)
	}

	if newHashErr != nil {
		return nil, nil, fmt.Errorf("new hash: %w", newHashErr)
	}

	return output[:cipher.KeySize], output[cipher.KeySize:], nil
}
