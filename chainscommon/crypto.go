package chainscommon

import (
	"fmt"
	"hash"
	"io"

	"github.com/lyreware/go-ratchet/keys"
	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

const cryptoMessageCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX

var (
	cryptoMessageCipherKDFSalt = make([]byte, cryptoMessageCipherKDFOutputLen)
	cryptoMessageCipherKDFInfo = []byte("message cipher")
)

func DeriveMessageCipherKeyAndNonce(messageKey keys.Message) (key []byte, nonce []byte, err error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash

		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	kdf := hkdf.New(
		getHasher,
		messageKey.Bytes,
		cryptoMessageCipherKDFSalt,
		cryptoMessageCipherKDFInfo,
	)
	kdfOutput := make([]byte, cryptoMessageCipherKDFOutputLen)

	_, err = io.ReadFull(kdf, kdfOutput)
	if err != nil {
		return nil, nil, fmt.Errorf("KDF: %w", err)
	}

	if newHashErr != nil {
		return nil, nil, fmt.Errorf("new hash: %w", newHashErr)
	}

	key = kdfOutput[:cipher.KeySize]
	nonce = kdfOutput[cipher.KeySize:]

	return key, nonce, nil
}
