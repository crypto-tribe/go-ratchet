package receivingchain

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-ratchet/messagechainscommon"
)

type Crypto interface {
	AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error)
	DecryptHeader(key keys.Header, encryptedHeader []byte) (header.Header, error)
	DecryptMessage(key keys.Message, encryptedData, auth []byte) ([]byte, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	return defaultCrypto{}
}

func (c defaultCrypto) AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
	var newHashErr error

	getHasher := func() hash.Hash {
		var hasher hash.Hash
		hasher, newHashErr = blake2b.New512(nil)

		return hasher
	}

	mac := hmac.New(getHasher, masterKey.Bytes)

	const masterKeyByte = 0x02
	if _, err := mac.Write([]byte{masterKeyByte}); err != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("write %d byte to MAC: %w", masterKeyByte, err)
	}

	newMasterKey := keys.MessageMaster{Bytes: mac.Sum(nil)}
	mac.Reset()

	const messageKeyByte = 0x01
	if _, err := mac.Write([]byte{messageKeyByte}); err != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("write %d byte to MAC: %w", messageKeyByte, err)
	}

	messageKey := keys.Message{Bytes: mac.Sum(nil)}

	if newHashErr != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("new hash: %w", newHashErr)
	}

	return newMasterKey, messageKey, nil
}

func (c defaultCrypto) DecryptHeader(key keys.Header, encryptedHeader []byte) (header.Header, error) {
	if len(encryptedHeader) <= cipher.NonceSizeX {
		return header.Header{}, fmt.Errorf("encrpted header too short, expected at least %d bytes", cipher.NonceSizeX+1)
	}

	decryptedHeaderBytes, err := c.decrypt(
		key.Bytes, encryptedHeader[:cipher.NonceSizeX], encryptedHeader[cipher.NonceSizeX:], nil)
	if err != nil {
		return header.Header{}, err
	}

	decryptedHeader, err := header.Decode(decryptedHeaderBytes)
	if err != nil {
		return header.Header{}, fmt.Errorf("decode decrypted header: %w", err)
	}

	return decryptedHeader, nil
}

func (c defaultCrypto) DecryptMessage(key keys.Message, encryptedData, auth []byte) ([]byte, error) {
	cipherKey, nonce, err := messagechainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	return c.decrypt(cipherKey, nonce, encryptedData, auth)
}

func (c defaultCrypto) decrypt(key, nonce, encryptedData, auth []byte) ([]byte, error) {
	cipher, err := cipher.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	data, err := cipher.Open(nil, nonce, encryptedData, auth)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return data, nil
}
