package receivingchain

import (
	"crypto/hmac"
	"errors"
	"hash"

	"github.com/lyreware/go-ratchet/chainscommon"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
)

// Crypto is a crypto for the receiving chain.
type Crypto interface {
	AdvanceChain(masterKey keys.Master) (keys.Master, keys.Message, error)
	DecryptHeader(key keys.Header, encryptedHeader []byte) (header.Header, error)
	DecryptMessage(key keys.Message, encryptedMessage, auth []byte) ([]byte, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() defaultCrypto {
	crypto := defaultCrypto{}

	return crypto
}

func (defaultCrypto) AdvanceChain(masterKey keys.Master) (keys.Master, keys.Message, error) {
	var newHasherErr error

	getHasher := func() hash.Hash {
		hasher, err := blake2b.New512(nil)
		newHasherErr = err

		return hasher
	}

	mac := hmac.New(getHasher, masterKey.Bytes)

	const masterKeyByte = 0x02

	_, err := mac.Write([]byte{masterKeyByte})
	if err != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrWriteMasterKeyByteToMAC, err)
	}

	newMasterKey := keys.Master{
		Bytes: mac.Sum(nil),
	}
	mac.Reset()

	const messageKeyByte = 0x01

	_, err = mac.Write([]byte{messageKeyByte})
	if err != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrWriteMessageKeyByteToMAC, err)
	}

	messageKey := keys.Message{
		Bytes: mac.Sum(nil),
	}

	if newHasherErr != nil {
		return keys.Master{}, keys.Message{}, errors.Join(ErrNewHasher, newHasherErr)
	}

	return newMasterKey, messageKey, nil
}

func (c defaultCrypto) DecryptHeader(
	key keys.Header,
	encryptedHeader []byte,
) (header.Header, error) {
	if len(encryptedHeader) <= cipher.NonceSizeX {
		return header.Header{}, ErrNotEnoughEncryptedHeaderBytes
	}

	decryptedHeaderBytes, err := c.decrypt(
		key.Bytes,
		encryptedHeader[:cipher.NonceSizeX],
		encryptedHeader[cipher.NonceSizeX:],
		nil,
	)
	if err != nil {
		return header.Header{}, err
	}

	decryptedHeader, err := header.Decode(decryptedHeaderBytes)
	if err != nil {
		return header.Header{}, errors.Join(ErrDecodeHeader, err)
	}

	return decryptedHeader, nil
}

func (c defaultCrypto) DecryptMessage(
	key keys.Message,
	encryptedMessage []byte,
	auth []byte,
) ([]byte, error) {
	cipherKey, nonce, err := chainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return nil, errors.Join(ErrDeriveMessageCipherKeyAndNonce, err)
	}

	decryptedMessage, err := c.decrypt(cipherKey, nonce, encryptedMessage, auth)
	if err != nil {
		return nil, errors.Join(ErrDecrypt, err)
	}

	return decryptedMessage, nil
}

func (defaultCrypto) decrypt(key, nonce, encryptedData, auth []byte) ([]byte, error) {
	cipherX, err := cipher.NewX(key)
	if err != nil {
		return nil, errors.Join(ErrNewCipher, err)
	}

	decryptedData, err := cipherX.Open(nil, nonce, encryptedData, auth)
	if err != nil {
		return nil, errors.Join(ErrOpenCipher, err)
	}

	return decryptedData, nil
}
