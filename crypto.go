package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"

	"github.com/lyreware/go-ratchet/keys"
)

type Crypto interface {
	ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error)
	GenerateKeyPair() (keys.Private, keys.Public, error)
}

type defaultCrypto struct {
	curve ecdh.Curve
}

func newDefaultCrypto() (crypto defaultCrypto) {
	crypto = defaultCrypto{
		curve: ecdh.X25519(),
	}

	return crypto
}

func (c defaultCrypto) ComputeSharedKey(
	privateKey keys.Private,
	publicKey keys.Public,
) (sharedKey keys.Shared, err error) {
	foreignPrivateKey, err := c.curve.NewPrivateKey(privateKey.Bytes)
	if err != nil {
		return sharedKey, fmt.Errorf("map to foreign private key: %w", err)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes)
	if err != nil {
		return sharedKey, fmt.Errorf("map to foreign public key: %w", err)
	}

	sharedKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return sharedKey, fmt.Errorf("Diffie-Hellman: %w", err)
	}

	sharedKey = keys.Shared{
		Bytes: sharedKeyBytes,
	}

	return sharedKey, err
}

func (c defaultCrypto) GenerateKeyPair() (privateKey keys.Private, publicKey keys.Public, err error) {
	foreignPrivateKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return privateKey, publicKey, err
	}

	privateKey = keys.Private{
		Bytes: foreignPrivateKey.Bytes(),
	}

	publicKey = keys.Public{
		Bytes: foreignPrivateKey.PublicKey().Bytes(),
	}

	return privateKey, publicKey, err
}
