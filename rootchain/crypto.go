package rootchain

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/platform-inf/go-ratchet/keys"
)

const defaultCryptoKDFOutputLen = 3 * 32

var defaultCryptoKDFInfo = []byte("advance root chain")

type Crypto interface {
	AdvanceChain(rootKey keys.Root, sharedKey keys.Shared) (keys.Root, keys.MessageMaster, keys.Header, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	return defaultCrypto{}
}

func (crypto defaultCrypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (keys.Root, keys.MessageMaster, keys.Header, error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash
		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	kdf := hkdf.New(getHasher, sharedKey.Bytes, rootKey.Bytes, defaultCryptoKDFInfo)

	output := make([]byte, defaultCryptoKDFOutputLen)
	if _, err := io.ReadFull(kdf, output); err != nil {
		return keys.Root{}, keys.MessageMaster{}, keys.Header{}, fmt.Errorf("KDF: %w", err)
	}

	if newHashErr != nil {
		return keys.Root{}, keys.MessageMaster{}, keys.Header{}, fmt.Errorf("new hash: %w", newHashErr)
	}

	newRootKey := keys.Root{Bytes: output[:32]}
	messageMasterKey := keys.MessageMaster{Bytes: output[32:64]}
	nextHeaderKey := keys.Header{Bytes: output[64:]}

	return newRootKey, messageMasterKey, nextHeaderKey, nil
}
