package receivingchain

import (
	"errors"
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-utils"
)

// Ratchet receiving chain.
//
// Please note that this structure may corrupt its state in case of errors. Therefore, clone the data at the top level
// and replace the current data with it if there are no errors.
type Chain struct {
	masterKey         *keys.MessageMaster
	headerKey         *keys.Header
	nextHeaderKey     keys.Header
	nextMessageNumber uint64
	cfg               config
}

func New(
	masterKey *keys.MessageMaster,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	options ...Option,
) (Chain, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	chain := Chain{
		masterKey:         masterKey,
		headerKey:         headerKey,
		nextHeaderKey:     nextHeaderKey,
		nextMessageNumber: nextMessageNumber,
		cfg:               cfg,
	}

	return chain, nil
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
) ([]byte, error) {
	decryptedData, err := ch.decryptWithSkippedKeys(encryptedHeader, encryptedData, auth)
	if err == nil {
		return decryptedData, nil
	}

	err = fmt.Errorf("decrypt with skipped keys: %w", err)

	if handleErr := ch.handleEncryptedHeader(encryptedHeader, ratchet); handleErr != nil {
		return nil, errors.Join(err, fmt.Errorf("handle encrypted header: %w", handleErr))
	}

	messageKey, advanceErr := ch.advance()
	if advanceErr != nil {
		return nil, errors.Join(err, fmt.Errorf("advance chain: %w", err))
	}

	auth = utils.ConcatByteSlices(encryptedHeader, auth)

	decryptedData, decryptErr := ch.cfg.crypto.DecryptMessage(messageKey, encryptedData, auth)
	if decryptErr != nil {
		return nil, errors.Join(err, fmt.Errorf("%w: decrypt message: %w", errlist.ErrCrypto, decryptErr))
	}

	// Note that here it is ok to ignore an error when decrypting with skipped keys if decryption with the next message key
	// succeeds.
	return decryptedData, nil
}

func (ch *Chain) Upgrade(masterKey keys.MessageMaster, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = &ch.nextHeaderKey
	ch.nextHeaderKey = nextHeaderKey
	ch.nextMessageNumber = 0
}

func (ch *Chain) advance() (keys.Message, error) {
	if ch.masterKey == nil {
		return keys.Message{}, fmt.Errorf("%w: master key is nil", errlist.ErrInvalidValue)
	}

	newMasterKey, messageKey, err := ch.cfg.crypto.AdvanceChain(*ch.masterKey)
	if err != nil {
		return keys.Message{}, fmt.Errorf("%w: advance via crypto: %w", errlist.ErrCrypto, err)
	}

	ch.masterKey = &newMasterKey
	ch.nextMessageNumber++

	return messageKey, nil
}

// decryptHeaderWithCurrentOrNextKeys must decrypt passed encrypted header with current or next header key.
//
// Note that ratchet is needed if header decrypted with next header key.
func (ch *Chain) decryptHeaderWithCurrentOrNextKey(
	encryptedHeader []byte,
) (decryptedHeader header.Header, needRatchet bool, err error) {
	if ch.headerKey != nil {
		header, decryptErr := ch.cfg.crypto.DecryptHeader(*ch.headerKey, encryptedHeader)
		if decryptErr == nil {
			return header, false, nil
		}

		err = errors.Join(err, fmt.Errorf("%w: decrypt header with current key: %w", errlist.ErrCrypto, decryptErr))
	}

	decryptedHeader, decryptErr := ch.cfg.crypto.DecryptHeader(ch.nextHeaderKey, encryptedHeader)
	if decryptErr != nil {
		err = errors.Join(err, fmt.Errorf("%w: decrypt header with next header key: %w", errlist.ErrCrypto, decryptErr))
		return header.Header{}, false, err
	}

	// Note that here it is ok to ignore an error when decrypting with the current key if decryption with the next key
	// succeeds.
	return decryptedHeader, true, nil
}

func (ch *Chain) decryptWithSkippedKeys(encryptedHeader, encryptedData, auth []byte) ([]byte, error) {
	iter, err := ch.cfg.skippedKeysStorage.GetIter()
	if err != nil {
		return nil, fmt.Errorf("%w: get iter: %w", errlist.ErrSkippedKeysStorage, err)
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
				return nil, fmt.Errorf("%w: decrypt message: %w", errlist.ErrCrypto, err)
			}

			if err := ch.cfg.skippedKeysStorage.Delete(headerKey, messageNumber); err != nil {
				return nil, fmt.Errorf("%w: delete: %w", errlist.ErrSkippedKeysStorage, err)
			}

			return decryptedData, nil
		}
	}

	return nil, errors.New("no keys to decrypt header and data")
}

func (ch *Chain) handleEncryptedHeader(encryptedHeader []byte, ratchet RatchetCallback) error {
	decryptedHeader, needRatchet, err := ch.decryptHeaderWithCurrentOrNextKey(encryptedHeader)
	if err != nil {
		return fmt.Errorf("decrypt header: %w", err)
	}

	if needRatchet {
		if err := ch.skipKeys(decryptedHeader.PreviousSendingChainMessagesCount); err != nil {
			return fmt.Errorf("skip %d keys: %w", decryptedHeader.PreviousSendingChainMessagesCount, err)
		}

		if err := ratchet(decryptedHeader.PublicKey); err != nil {
			return fmt.Errorf("ratchet: %w", err)
		}
	}

	if err := ch.skipKeys(decryptedHeader.MessageNumber); err != nil {
		return fmt.Errorf("skip %d message keys in upgraded chain: %w", decryptedHeader.MessageNumber, err)
	}

	return nil
}

func (ch *Chain) skipKeys(untilMessageNumber uint64) error {
	if untilMessageNumber < ch.nextMessageNumber {
		return fmt.Errorf("message number is small for the current chain, next message number is %d", ch.nextMessageNumber)
	}

	for messageNumber := ch.nextMessageNumber; messageNumber < untilMessageNumber; messageNumber++ {
		messageKey, err := ch.advance()
		if err != nil {
			return fmt.Errorf("advance chain: %w", err)
		}

		if ch.headerKey == nil {
			return fmt.Errorf("%w: header key is nil", errlist.ErrInvalidValue)
		}

		if err := ch.cfg.skippedKeysStorage.Add(*ch.headerKey, messageNumber, messageKey); err != nil {
			return fmt.Errorf("%w: add: %w", errlist.ErrSkippedKeysStorage, err)
		}
	}

	return nil
}

// RatchetCallback must perform ratchet and upgrade receiving chain.
type RatchetCallback func(remotePublicKey keys.Public) error
