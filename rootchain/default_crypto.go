package rootchain

import (
	"errors"
	"hash"
	"io"

	"github.com/platform-source/aegis/keys"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"
)

var defaultCryptoKDFInfo = []byte("advance root chain")

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	crypto := defaultCrypto{}

	return crypto
}

func (defaultCrypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	const kdfOutputKeySize = 32

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
	kdfOutput := make([]byte, 3*kdfOutputKeySize)

	_, err := io.ReadFull(kdf, kdfOutput)
	if err != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, errors.Join(ErrKDF, err)
	}

	if newHashErr != nil {
		return keys.Root{}, keys.Master{}, keys.Header{}, errors.Join(ErrNewHasher, newHashErr)
	}

	newRootKey := keys.Root{
		Bytes: kdfOutput[:kdfOutputKeySize],
	}

	masterKey := keys.Master{
		Bytes: kdfOutput[kdfOutputKeySize : 2*kdfOutputKeySize],
	}

	nextHeaderKey := keys.Header{
		Bytes: kdfOutput[2*kdfOutputKeySize : 3*kdfOutputKeySize],
	}

	return newRootKey, masterKey, nextHeaderKey, nil
}
