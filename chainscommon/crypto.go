package chainscommon

import (
	"errors"
	"hash"
	"io"

	"github.com/lyreware/go-ratchet/keys"
	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

const (
	messageCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX
)

var (
	messageCipherKDFSalt = make([]byte, messageCipherKDFOutputLen)
	messageCipherKDFInfo = []byte("message cipher")
)

// DeriveMessageCipherKeyAndNonce derives a new cipher key and cipher nonce to encrypt a message.
func DeriveMessageCipherKeyAndNonce(messageKey keys.Message) (key []byte, nonce []byte, err error) {
	var newHasherErr error

	kdf := hkdf.New(
		func() hash.Hash {
			hasher, err := blake2b.New512(nil)
			newHasherErr = err

			return hasher
		},
		messageKey.Bytes,
		messageCipherKDFSalt,
		messageCipherKDFInfo,
	)
	kdfOutput := make([]byte, messageCipherKDFOutputLen)

	_, err = io.ReadFull(kdf, kdfOutput)
	if err != nil {
		return nil, nil, errors.Join(ErrKDF, err)
	}

	if newHasherErr != nil {
		return nil, nil, errors.Join(ErrNewHasher, newHasherErr)
	}

	key = kdfOutput[:cipher.KeySize]
	nonce = kdfOutput[cipher.KeySize:]

	return key, nonce, nil
}
