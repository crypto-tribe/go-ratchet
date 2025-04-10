package sendingchain

import (
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-utils"
)

// Ratchet sending chain.
//
// Please note that this structure may corrupt its state in case of errors. Therefore, clone the data at the top level
// and replace the current data with it if there are no errors.
type Chain struct {
	masterKey                  *keys.MessageMaster
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	cfg                        config
}

func New(
	masterKey *keys.MessageMaster,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	previousChainMessagesCount uint64,
	options ...Option,
) (Chain, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	chain := Chain{
		masterKey:                  masterKey,
		headerKey:                  headerKey,
		nextHeaderKey:              nextHeaderKey,
		nextMessageNumber:          nextMessageNumber,
		previousChainMessagesCount: previousChainMessagesCount,
		cfg:                        cfg,
	}

	return chain, nil
}

func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()

	return ch
}

func (ch *Chain) Encrypt(header header.Header, data, auth []byte) ([]byte, []byte, error) {
	if ch.headerKey == nil {
		return nil, nil, fmt.Errorf("%w: header key is nil", errlist.ErrInvalidValue)
	}

	encryptedHeader, err := ch.cfg.crypto.EncryptHeader(*ch.headerKey, header)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: encrypt header: %w", errlist.ErrCrypto, err)
	}

	messageKey, err := ch.advance()
	if err != nil {
		return nil, nil, fmt.Errorf("advance chain: %w", err)
	}

	auth = utils.ConcatByteSlices(encryptedHeader, auth)

	encryptedData, err := ch.cfg.crypto.EncryptMessage(messageKey, data, auth)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: encrypt message: %w", errlist.ErrCrypto, err)
	}

	return encryptedHeader, encryptedData, nil
}

func (ch *Chain) PrepareHeader(publicKey keys.Public) header.Header {
	return header.Header{
		PublicKey:                         publicKey,
		MessageNumber:                     ch.nextMessageNumber,
		PreviousSendingChainMessagesCount: ch.previousChainMessagesCount,
	}
}

func (ch *Chain) Upgrade(masterKey keys.MessageMaster, nextHeaderKey keys.Header) {
	headerKey := ch.nextHeaderKey.Clone()

	ch.masterKey = &masterKey
	ch.headerKey = &headerKey
	ch.nextHeaderKey = nextHeaderKey
	ch.previousChainMessagesCount = ch.nextMessageNumber
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
