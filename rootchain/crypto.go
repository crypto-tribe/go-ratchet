package rootchain

import (
	"errors"
	"hash"
	"io"

	"github.com/lyreware/go-ratchet/keys"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"
)

const (
	defaultCryptoKDFOutputLen = 3 * 32
)

var defaultCryptoKDFInfo = []byte("advance root chain")

// Crypto is the crypto interface for the root chain.
type Crypto interface {
	AdvanceChain(
		rootKey keys.Root,
		sharedKey keys.Shared,
	) (keys.Root, keys.Master, keys.Header, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	crypto := defaultCrypto{}

	return crypto
}

func (defaultCrypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	var newHashErr error

	kdf := hkdf.New(
		func() hash.Hash {
			hasher, err := blake2b.New512(nil)
			newHashErr = err

			return hasher
		},
		sharedKey.Bytes,
		rootKey.Bytes,
		defaultCryptoKDFInfo,
	)
	kdfOutput := make([]byte, defaultCryptoKDFOutputLen)

	_, err := io.ReadFull(kdf, kdfOutput)
	if err != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, errors.Join(ErrKDF, err)
	}

	if newHashErr != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, errors.Join(ErrNewHasher, newHashErr)
	}

	newRootKey := keys.Root{
		Bytes: kdfOutput[:32],
	}

	masterKey := keys.Master{
		Bytes: kdfOutput[32:64],
	}

	nextHeaderKey := keys.Header{
		Bytes: kdfOutput[64:],
	}

	return newRootKey, masterKey, nextHeaderKey, nil
}
