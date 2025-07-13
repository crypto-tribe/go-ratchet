package rootchain

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/lyreware/go-ratchet/keys"
)

const defaultCryptoKDFOutputLen = 3 * 32

var defaultCryptoKDFInfo = []byte("advance root chain")

type Crypto interface {
	AdvanceChain(
		rootKey keys.Root,
		sharedKey keys.Shared,
	) (newRootKey keys.Root, masterKey keys.Master, nextHeaderKey keys.Header, err error)
}

type defaultCrypto struct{}

func newDefaultCrypto() (crypto defaultCrypto) {
	return crypto
}

func (crypto defaultCrypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (newRootKey keys.Root, masterKey keys.Master, nextHeaderKey keys.Header, err error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash

		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	kdf := hkdf.New(getHasher, sharedKey.Bytes, rootKey.Bytes, defaultCryptoKDFInfo)
	output := make([]byte, defaultCryptoKDFOutputLen)

	_, err = io.ReadFull(kdf, output)
	if err != nil {
		return newRootKey, masterKey, nextHeaderKey, fmt.Errorf("KDF: %w", err)
	}

	if newHashErr != nil {
		return newRootKey, masterKey, nextHeaderKey, fmt.Errorf("new hash: %w", newHashErr)
	}

	newRootKey = keys.Root{Bytes: output[:32]}
	masterKey = keys.Master{Bytes: output[32:64]}
	nextHeaderKey = keys.Header{Bytes: output[64:]}

	return newRootKey, masterKey, nextHeaderKey, err
}
