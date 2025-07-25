package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"

	"github.com/platform-source/aegis/keys"
)

type defaultCrypto struct {
	curve ecdh.Curve
}

func newDefaultCrypto() defaultCrypto {
	crypto := defaultCrypto{
		curve: ecdh.X25519(),
	}

	return crypto
}

func (c defaultCrypto) ComputeSharedKey(
	privateKey keys.Private,
	publicKey keys.Public,
) (keys.Shared, error) {
	foreignPrivateKey, err := c.curve.NewPrivateKey(privateKey.Bytes)
	if err != nil {
		return keys.Shared{}, errors.Join(ErrNewPrivateKey, err)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes)
	if err != nil {
		return keys.Shared{}, errors.Join(ErrNewPublicKey, err)
	}

	sharedKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return keys.Shared{}, errors.Join(ErrDiffieHellman, err)
	}

	sharedKey := keys.Shared{
		Bytes: sharedKeyBytes,
	}

	return sharedKey, nil
}

func (c defaultCrypto) GenerateKeyPair() (keys.Private, keys.Public, error) {
	foreignPrivateKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return keys.Private{}, keys.Public{}, errors.Join(ErrGeneratePrivateKey, err)
	}

	privateKey := keys.Private{
		Bytes: foreignPrivateKey.Bytes(),
	}

	publicKey := keys.Public{
		Bytes: foreignPrivateKey.PublicKey().Bytes(),
	}

	return privateKey, publicKey, nil
}
