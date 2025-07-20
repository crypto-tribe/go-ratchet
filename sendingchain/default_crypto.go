package sendingchain

import (
	"crypto/hmac"
	"crypto/rand"
	"errors"
	"hash"

	"github.com/crypto-tribe/go-ratchet/chainscommon"
	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
	"github.com/crypto-tribe/go-utils/slices"
	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
)

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	crypto := defaultCrypto{}

	return crypto
}

func (defaultCrypto) AdvanceChain(masterKey keys.Master) (keys.Master, keys.Message, error) {
	const (
		masterKeyByte  = 0x02
		messageKeyByte = 0x01
	)

	var newHashErr error

	mac := hmac.New(
		func() hash.Hash {
			hasher, err := blake2b.New512(nil)
			newHashErr = err

			return hasher
		},
		masterKey.Bytes,
	)

	_, err := mac.Write([]byte{masterKeyByte})
	if err != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrWriteMasterKeyByteToMAC, err)
	}

	newMasterKey := keys.Master{
		Bytes: mac.Sum(nil),
	}
	mac.Reset()

	_, err = mac.Write([]byte{messageKeyByte})
	if err != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrWriteMessageKeyByteToMAC, err)
	}

	messageKey := keys.Message{
		Bytes: mac.Sum(nil),
	}

	if newHashErr != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrNewHasher, newHashErr)
	}

	return newMasterKey, messageKey, nil
}

func (c defaultCrypto) EncryptHeader(key keys.Header, head header.Header) ([]byte, error) {
	var nonce [cipher.NonceSizeX]byte

	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, errors.Join(ErrGenerateNonce, err)
	}

	encryptedHeader, err := c.encrypt(key.Bytes, nonce[:], head.Encode(), nil)
	if err != nil {
		return nil, errors.Join(ErrEncrypt, err)
	}

	encryptedHeader = slices.ConcatBytes(nonce[:], encryptedHeader)

	return encryptedHeader, nil
}

func (c defaultCrypto) EncryptMessage(key keys.Message, message, auth []byte) ([]byte, error) {
	cipherKey, nonce, err := chainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return nil, errors.Join(ErrDeriveMessageCipherKeyAndNonce, err)
	}

	encryptedMessage, err := c.encrypt(cipherKey, nonce, message, auth)
	if err != nil {
		return nil, errors.Join(ErrEncrypt, err)
	}

	return encryptedMessage, nil
}

func (defaultCrypto) encrypt(key, nonce, data, auth []byte) ([]byte, error) {
	cipherX, err := cipher.NewX(key)
	if err != nil {
		return nil, errors.Join(ErrNewCipher, err)
	}

	encryptedData := cipherX.Seal(nil, nonce, data, auth)

	return encryptedData, nil
}
