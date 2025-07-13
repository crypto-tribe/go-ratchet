package receivingchain

import (
	"errors"
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-utils"
)

// Ratchet receiving chain.
//
// Please note that this structure may corrupt its state in case of errors. Therefore, clone the data at the top level
// and replace the current data with it if there are no errors.
type Chain struct {
	masterKey         *keys.Master
	headerKey         *keys.Header
	nextHeaderKey     keys.Header
	nextMessageNumber uint64
	cfg               config
}

func New(
	masterKey *keys.Master,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	options ...Option,
) (chain Chain, err error) {
	chain = Chain{
		masterKey:         masterKey,
		headerKey:         headerKey,
		nextHeaderKey:     nextHeaderKey,
		nextMessageNumber: nextMessageNumber,
	}

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return chain, fmt.Errorf("new config: %w", err)
	}

	return chain, err
}

func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()
	ch.cfg = ch.cfg.clone()

	return ch
}

func (ch *Chain) Decrypt(
	encryptedHeader []byte,
	encryptedData []byte,
	auth []byte,
	ratchet RatchetCallback,
) (decryptedData []byte, err error) {
	decryptedData, skippedKeysErr := ch.decryptWithSkippedKeys(encryptedHeader, encryptedData, auth)
	if err == nil {
		return decryptedData, err
	}

	skippedKeysErr = fmt.Errorf("decrypt with skipped keys: %w", skippedKeysErr)

	err = ch.handleEncryptedHeader(encryptedHeader, ratchet)
	if err != nil {
		err = errors.Join(skippedKeysErr, fmt.Errorf("handle encrypted header: %w", err))
		return decryptedData, err
	}

	messageKey, err := ch.advance()
	if err != nil {
		err = errors.Join(skippedKeysErr, fmt.Errorf("advance chain: %w", err))
		return decryptedData, err
	}

	auth = utils.ConcatByteSlices(encryptedHeader, auth)

	decryptedData, err = ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
	if err != nil {
		err = errors.Join(skippedKeysErr, fmt.Errorf("%w: decrypt message: %w", errlist.ErrCrypto, err))
		return decryptedData, err
	}

	// Note that here it is ok to ignore an error when decrypting with skipped keys if decryption with the next message key
	// succeeds.
	return decryptedData, err
}

func (ch *Chain) Upgrade(masterKey keys.Master, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.nextMessageNumber = 0
}

func (ch *Chain) advance() (messageKey keys.Message, err error) {
	if ch.masterKey == nil {
		return messageKey, fmt.Errorf("%w: master key is nil", errlist.ErrInvalidValue)
	}

	var newMasterKey keys.Master

	newMasterKey, messageKey, err = ch.cfg.crypto.AdvanceChain(*ch.masterKey)
	if err != nil {
		return messageKey, fmt.Errorf("%w: advance via crypto: %w", errlist.ErrCrypto, err)
	}

	ch.masterKey = &newMasterKey
	ch.nextMessageNumber++

	return messageKey, err
}

// decryptHeaderWithCurrentOrNextKeys must decrypt passed encrypted header with current or next header key.
//
// Note that ratchet is needed if header decrypted with next header key.
func (ch *Chain) decryptHeaderWithCurrentOrNextKey(
	encryptedHeader []byte,
) (decryptedHeader header.Header, needRatchet bool, err error) {
	var currentKeyErr error

	if ch.headerKey != nil {
		decryptedHeader, currentKeyErr = ch.cfg.crypto.DecryptHeader(*ch.headerKey, encryptedHeader)
		if currentKeyErr == nil {
			return decryptedHeader, false, err
		}

		currentKeyErr = fmt.Errorf("%w: decrypt header with current key: %w", errlist.ErrCrypto, currentKeyErr)
	}

	decryptedHeader, err = ch.cfg.crypto.DecryptHeader(ch.nextHeaderKey, encryptedHeader)
	if err != nil {
		err = errors.Join(
			currentKeyErr,
			fmt.Errorf("%w: decrypt header with next header key: %w", errlist.ErrCrypto, err),
		)
		return decryptedHeader, false, err
	}

	// Note that here it is ok to ignore an error when decrypting with the current key if decryption with the next key
	// succeeds.
	return decryptedHeader, true, err
}

func (ch *Chain) decryptWithSkippedKeys(encryptedHeader, encryptedData, auth []byte) (decryptedData []byte, err error) {
	iter, err := ch.cfg.skippedKeysStorage.GetIter()
	if err != nil {
		return decryptedData, fmt.Errorf("%w: get iter: %w", errlist.ErrSkippedKeysStorage, err)
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

			decryptedData, err = ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
			if err != nil {
				return decryptedData, fmt.Errorf("%w: decrypt message: %w", errlist.ErrCrypto, err)
			}

			err = ch.cfg.skippedKeysStorage.Delete(headerKey, messageNumber)
			if err != nil {
				return decryptedData, fmt.Errorf("%w: delete: %w", errlist.ErrSkippedKeysStorage, err)
			}

			return decryptedData, err
		}
	}

	return decryptedData, errors.New("no keys to decrypt header and data")
}

func (ch *Chain) handleEncryptedHeader(encryptedHeader []byte, ratchet RatchetCallback) (err error) {
	decryptedHeader, needRatchet, err := ch.decryptHeaderWithCurrentOrNextKey(encryptedHeader)
	if err != nil {
		return fmt.Errorf("decrypt header: %w", err)
	}

	if needRatchet {
		err = ch.skipKeys(decryptedHeader.PreviousSendingChainMessagesCount)
		if err != nil {
			return fmt.Errorf("skip %d keys: %w", decryptedHeader.PreviousSendingChainMessagesCount, err)
		}

		err = ratchet(decryptedHeader.PublicKey)
		if err != nil {
			return fmt.Errorf("ratchet: %w", err)
		}
	}

	err = ch.skipKeys(decryptedHeader.MessageNumber)
	if err != nil {
		return fmt.Errorf("skip %d message keys in upgraded chain: %w", decryptedHeader.MessageNumber, err)
	}

	return err
}

func (ch *Chain) skipKeys(untilMessageNumber uint64) (err error) {
	if untilMessageNumber < ch.nextMessageNumber {
		return fmt.Errorf(
			"message number is small for the current chain, next message number is %d",
			ch.nextMessageNumber,
		)
	}

	for messageNumber := ch.nextMessageNumber; messageNumber < untilMessageNumber; messageNumber++ {
		messageKey, err := ch.advance()
		if err != nil {
			return fmt.Errorf("advance chain: %w", err)
		}

		if ch.headerKey == nil {
			return fmt.Errorf("%w: header key is nil", errlist.ErrInvalidValue)
		}

		err = ch.cfg.skippedKeysStorage.Add(*ch.headerKey, messageNumber, messageKey)
		if err != nil {
			return fmt.Errorf("%w: add: %w", errlist.ErrSkippedKeysStorage, err)
		}
	}

	return err
}

// RatchetCallback must perform ratchet and upgrade receiving chain.
type RatchetCallback func(remotePublicKey keys.Public) error
