package sendingchain

import (
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-ratchet/messagechainscommon"
	"github.com/lyreware/go-utils"
)

type Crypto interface {
	AdvanceChain(masterKey keys.Master) (newMasterKey keys.Master, messageKey keys.Message, err error)
	EncryptHeader(key keys.Header, header header.Header) (encryptedHeader []byte, err error)
	EncryptMessage(key keys.Message, message, auth []byte) (encryptedMessage []byte, err error)
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

	newMasterKey = keys.Master{
		Bytes: mac.Sum(nil),
	}
	mac.Reset()

	const messageKeyByte = 0x01

	_, err = mac.Write([]byte{messageKeyByte})
	if err != nil {
		return newMasterKey, messageKey, fmt.Errorf("write %d byte to MAC: %w", messageKeyByte, err)
	}

	messageKey = keys.Message{
		Bytes: mac.Sum(nil),
	}

	if newHashErr != nil {
		return newMasterKey, messageKey, fmt.Errorf("new hash: %w", newHashErr)
	}

	return newMasterKey, messageKey, err
}

func (c defaultCrypto) EncryptHeader(key keys.Header, header header.Header) (encryptedHeader []byte, err error) {
	var nonce [cipher.NonceSizeX]byte

	_, err = rand.Read(nonce[:])
	if err != nil {
		return bytes, fmt.Errorf("generate random nonce: %w", err)
	}

	encryptedHeader, err := c.encrypt(key.Bytes, nonce[:], header.Encode(), nil)
	if err != nil {
		return bytes, fmt.Errorf("encrypt: %w", err)
	}

	encryptedHeader = utils.ConcatByteSlices(nonce[:], encryptedHeader)

	return encryptedHeader, err
}

func (c defaultCrypto) EncryptMessage(key keys.Message, message, auth []byte) (encryptedMessage []byte, err error) {
	cipherKey, nonce, err := messagechainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return bytes, fmt.Errorf("derive key and nonce: %w", err)
	}

	encryptedMessage, err = c.encrypt(cipherKey, nonce, message, auth)

	return encryptedMessage, err
}

func (c defaultCrypto) encrypt(key, nonce, data, auth []byte) (encryptedData []byte, err error) {
	cipher, err := cipher.NewX(key)
	if err != nil {
		return bytes, fmt.Errorf("new cipher: %w", err)
	}

	encryptedData = cipher.Seal(nil, nonce, data, auth)

	return encryptedData, err
}
