package rootchain

import (
	"fmt"
	"hash"
	"io"

	"github.com/lyreware/go-ratchet/keys"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"
)

const defaultCryptoKDFOutputLen = 3 * 32

var defaultCryptoKDFInfo = []byte("advance root chain")

// Crypto is the crypto interface for the root chain.
type Crypto interface {
	AdvanceChain(
		rootKey keys.Root,
		sharedKey keys.Shared,
	) (keys.Root, keys.Master, keys.Header, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() (crypto defaultCrypto) {
	return crypto
}

func (defaultCrypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash

		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	kdf := hkdf.New(getHasher, sharedKey.Bytes, rootKey.Bytes, defaultCryptoKDFInfo)
	output := make([]byte, defaultCryptoKDFOutputLen)

	_, err := io.ReadFull(kdf, output)
	if err != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, fmt.Errorf("KDF: %w", err)
	}

	if newHashErr != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, fmt.Errorf("new hash: %w", newHashErr)
	}

	newRootKey := keys.Root{
		Bytes: output[:32],
	}

	masterKey := keys.Master{
		Bytes: output[32:64],
	}

	nextHeaderKey := keys.Header{
		Bytes: output[64:],
	}

	return newRootKey, masterKey, nextHeaderKey, nil
}
