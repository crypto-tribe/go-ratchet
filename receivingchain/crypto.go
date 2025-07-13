package receivingchain

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-ratchet/messagechainscommon"
)

type Crypto interface {
	AdvanceChain(masterKey keys.Master) (newMasterKey keys.Master, messageKey keys.Message, err error)
	DecryptHeader(key keys.Header, encryptedHeader []byte) (header.Header, error)
	DecryptMessage(key keys.Message, encryptedMessage, auth []byte) (decryptedMessage []byte, err error)
}

type defaultCrypto struct{}

func newDefaultCrypto() (crypto defaultCrypto) {
	return crypto
}

func (c defaultCrypto) AdvanceChain(
	masterKey keys.Master,
) (newMasterKey keys.Master, messageKey keys.Message, err error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash

		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	mac := hmac.New(getHasher, masterKey.Bytes)

	const masterKeyByte = 0x02

	_, err = mac.Write([]byte{masterKeyByte})
	if err != nil {
		return newMasterKey, messageKey, fmt.Errorf("write %d byte to MAC: %w", masterKeyByte, err)
	}

	newMasterKey.Bytes = mac.Sum(nil)
	mac.Reset()

	const messageKeyByte = 0x01

	_, err = mac.Write([]byte{messageKeyByte})
	if err != nil {
		return newMasterKey, messageKey, fmt.Errorf("write %d byte to MAC: %w", messageKeyByte, err)
	}

	messageKey.Bytes = mac.Sum(nil)

	if newHashErr != nil {
		return newMasterKey, messageKey, fmt.Errorf("new hash: %w", newHashErr)
	}

	return newMasterKey, messageKey, err
}

func (c defaultCrypto) DecryptHeader(
	key keys.Header,
	encryptedHeader []byte,
) (decryptedHeader header.Header, err error) {
	if len(encryptedHeader) <= cipher.NonceSizeX {
		return decryptedHeader, fmt.Errorf("encrpted header too short, expected at least %d bytes", cipher.NonceSizeX+1)
	}

	decryptedHeaderBytes, err := c.decrypt(
		key.Bytes,
		encryptedHeader[:cipher.NonceSizeX],
		encryptedHeader[cipher.NonceSizeX:],
		nil,
	)
	if err != nil {
		return decryptedHeader, fmt.Errorf("decrypt: %w", err)
	}

	decryptedHeader, err = header.Decode(decryptedHeaderBytes)
	if err != nil {
		return decryptedHeader, fmt.Errorf("decode decrypted header: %w", err)
	}

	return decryptedHeader, err
}

func (c defaultCrypto) DecryptMessage(
	key keys.Message,
	encryptedMessage []byte,
	auth []byte,
) (decryptedMessage []byte, err error) {
	cipherKey, nonce, err := messagechainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return decryptedData, fmt.Errorf("derive key and nonce: %w", err)
	}

	decryptedData, err = c.decrypt(cipherKey, nonce, encryptedData, auth)
	if err != nil {
		return decryptedData, fmt.Errorf("decrypt: %w", err)
	}

	return decryptedData, err
}

func (c defaultCrypto) decrypt(key, nonce, encryptedData, auth []byte) (decryptedData []byte, err error) {
	cipher, err := cipher.NewX(key)
	if err != nil {
		return decryptedData, fmt.Errorf("new cipher: %w", err)
	}

	decryptedData, err = cipher.Open(nil, nonce, encryptedData, auth)
	if err != nil {
		return decryptedData, fmt.Errorf("decrypt: %w", err)
	}

	return decryptedData, err
}
