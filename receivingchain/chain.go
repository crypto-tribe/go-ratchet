package receivingchain

import (
	"errors"

	"github.com/platform-source/aegis/header"
	"github.com/platform-source/aegis/keys"
	"github.com/platform-source/tools/slices"
)

// Chain is the ratchet receiving chain.
//
// Please note that this structure may corrupt its state in case of errors.
// Therefore, clone the data at the top level and replace the current data
// with it if there are no errors.
type Chain struct {
	masterKey         *keys.Master
	headerKey         *keys.Header
	nextHeaderKey     keys.Header
	nextMessageNumber uint64
	cfg               config
}

// New creates a new receiving chain.
func New(
	masterKey *keys.Master,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	options ...Option,
) (Chain, error) {
	chain := Chain{
		masterKey:         masterKey,
		headerKey:         headerKey,
		nextHeaderKey:     nextHeaderKey,
		nextMessageNumber: nextMessageNumber,
	}

	var err error

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return Chain{}, errors.Join(ErrNewConfig, err)
	}

	return chain, nil
}

// Clone clones receiving chain.
func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()
	ch.cfg = ch.cfg.clone()

	return ch
}

// Decrypt decrypts passed encrypted header and encrypted data and authenticates
// them with auth. Also calls ratchet callback if ratchet is needed.
func (ch *Chain) Decrypt(
	encryptedHeader []byte,
	encryptedData []byte,
	auth []byte,
	ratchet RatchetCallback,
) ([]byte, error) {
	decryptedData, skippedKeysErr := ch.decryptWithSkippedKeys(encryptedHeader, encryptedData, auth)
	if skippedKeysErr == nil {
		return decryptedData, nil
	}

	skippedKeysErr = errors.Join(ErrDecryptWithSkippedKeys, skippedKeysErr)

	err := ch.handleEncryptedHeader(encryptedHeader, ratchet)
	if err != nil {
		return nil, errors.Join(skippedKeysErr, ErrHandleEncryptedHeader, err)
	}

	messageKey, err := ch.advance()
	if err != nil {
		return nil, errors.Join(skippedKeysErr, ErrAdvanceChain, err)
	}

	auth = slices.ConcatBytes(encryptedHeader, auth)

	decryptedData, err = ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
	if err != nil {
		return nil, errors.Join(skippedKeysErr, ErrDecryptMessage, err)
	}

	// Note that here it is ok to ignore an error when decrypting with skipped keys
	// if decryption with the next message key succeeds.
	return decryptedData, nil
}

// Upgrade upgrades receiving chain with new starting values.
func (ch *Chain) Upgrade(masterKey keys.Master, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.nextMessageNumber = 0
}

func (ch *Chain) advance() (keys.Message, error) {
	if ch.masterKey == nil {
		return keys.Message{}, ErrMasterKeyIsNil
	}

	newMasterKey, messageKey, err := ch.cfg.crypto.AdvanceChain(*ch.masterKey)
	if err != nil {
		return keys.Message{}, errors.Join(ErrCryptoAdvanceChain, err)
	}

	ch.masterKey = &newMasterKey
	ch.nextMessageNumber++

	return messageKey, nil
}

// decryptHeaderWithCurrentOrNextKeys must decrypt passed encrypted header with
// current or next header key.
//
// Note that ratchet is needed if header decrypted with next header key.
func (ch *Chain) decryptHeaderWithCurrentOrNextKey(
	encryptedHeader []byte,
) (decryptedHeader header.Header, needRatchet bool, err error) {
	var currentKeyErr error

	if ch.headerKey != nil {
		decryptedHeader, currentKeyErr = ch.cfg.crypto.DecryptHeader(*ch.headerKey, encryptedHeader)
		if currentKeyErr == nil {
			return decryptedHeader, false, nil
		}

		currentKeyErr = errors.Join(ErrDecryptHeaderWithCurrentKey, currentKeyErr)
	}

	decryptedHeader, err = ch.cfg.crypto.DecryptHeader(ch.nextHeaderKey, encryptedHeader)
	if err != nil {
		return header.Header{}, false, errors.Join(currentKeyErr, ErrDecryptHeaderWithNextKey, err)
	}

	// Note that here it is ok to ignore an error when decrypting with the current
	// key if decryption with the next key succeeds.
	return decryptedHeader, true, nil
}

func (ch *Chain) decryptWithSkippedKeys(
	encryptedHeader, encryptedData, auth []byte,
) ([]byte, error) {
	iter, err := ch.cfg.skippedKeysStorage.GetIter()
	if err != nil {
		return nil, errors.Join(ErrGetSkippedKeysStorageIter, err)
	}

	for headerKey, messageNumberKeys := range iter {
		decryptedHeader, err := ch.cfg.crypto.DecryptHeader(headerKey, encryptedHeader)
		if err != nil {
			continue
		}

		for messageNumber, messageKey := range messageNumberKeys {
			if messageNumber != decryptedHeader.MessageNumber {
				continue
			}

			decryptedData, err := ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
			if err != nil {
				return nil, errors.Join(ErrDecryptMessage, err)
			}

			err = ch.cfg.skippedKeysStorage.Delete(headerKey, messageNumber)
			if err != nil {
				return nil, errors.Join(ErrDeleteSkippedKeys, err)
			}

			return decryptedData, nil
		}
	}

	return nil, ErrSkippedKeysNotFound
}

func (ch *Chain) handleEncryptedHeader(encryptedHeader []byte, ratchet RatchetCallback) error {
	decryptedHeader, needRatchet, err := ch.decryptHeaderWithCurrentOrNextKey(encryptedHeader)
	if err != nil {
		return errors.Join(ErrDecryptHeaderWithCurrentOrNextKey, err)
	}

	if needRatchet {
		err = ch.skipKeys(decryptedHeader.PreviousSendingChainMessagesCount)
		if err != nil {
			return errors.Join(ErrSkipPreviousChainKeys, err)
		}

		err = ratchet(decryptedHeader.PublicKey)
		if err != nil {
			return errors.Join(ErrRatchet, err)
		}
	}

	err = ch.skipKeys(decryptedHeader.MessageNumber)
	if err != nil {
		return errors.Join(ErrSkipCurrentChainKeys, err)
	}

	return nil
}

func (ch *Chain) skipKeys(untilMessageNumber uint64) error {
	for messageNumber := ch.nextMessageNumber; messageNumber < untilMessageNumber; messageNumber++ {
		messageKey, err := ch.advance()
		if err != nil {
			return errors.Join(ErrAdvanceChain, err)
		}

		if ch.headerKey == nil {
			return ErrHeaderKeyIsNil
		}

		err = ch.cfg.skippedKeysStorage.Add(*ch.headerKey, messageNumber, messageKey)
		if err != nil {
			return errors.Join(ErrAddSkippedKey, err)
		}
	}

	return nil
}

// RatchetCallback must perform ratchet and upgrade receiving chain.
type RatchetCallback func(remotePublicKey keys.Public) error
