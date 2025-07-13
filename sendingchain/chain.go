package sendingchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-utils/convert"
)

// Ratchet sending chain.
//
// Please note that this structure may corrupt its state in case of errors. Therefore, clone the data at the top level
// and replace the current data with it if there are no errors.
type Chain struct {
	masterKey                  *keys.Master
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	cfg                        config
}

func New(
	masterKey *keys.Master,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	previousChainMessagesCount uint64,
	options ...Option,
) (chain Chain, err error) {
	chain = Chain{
		masterKey:                  masterKey,
		headerKey:                  headerKey,
		nextHeaderKey:              nextHeaderKey,
		nextMessageNumber:          nextMessageNumber,
		previousChainMessagesCount: previousChainMessagesCount,
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

	return ch
}

func (ch *Chain) Encrypt(
	header header.Header,
	data []byte,
	auth []byte,
) (encryptedHeader []byte, encryptedData []byte, err error) {
	if ch.headerKey == nil {
		return encryptedHeader, encryptedData, fmt.Errorf("%w: header key is nil", errlist.ErrInvalidValue)
	}

	encryptedHeader, err = ch.cfg.crypto.EncryptHeader(*ch.headerKey, header)
	if err != nil {
		return encryptedHeader, encryptedData, fmt.Errorf("%w: encrypt header: %w", errlist.ErrCrypto, err)
	}

	messageKey, err := ch.advance()
	if err != nil {
		return encryptedHeader, encryptedData, fmt.Errorf("advance chain: %w", err)
	}

	auth = utils.ConcatByteSlices(encryptedHeader, auth)

	encryptedData, err = ch.cfg.crypto.EncryptMessage(messageKey, data, auth)
	if err != nil {
		return encryptedHeader, encryptedData, fmt.Errorf("%w: encrypt message: %w", errlist.ErrCrypto, err)
	}

	return encryptedHeader, encryptedData, err
}

func (ch *Chain) PrepareHeader(publicKey keys.Public) (h header.Header) {
	h = header.Header{
		PublicKey:                         publicKey,
		MessageNumber:                     ch.nextMessageNumber,
		PreviousSendingChainMessagesCount: ch.previousChainMessagesCount,
	}

	return h
}

func (ch *Chain) Upgrade(masterKey keys.Master, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = convert.ToPtr(ch.nextHeaderKey.Clone())
	ch.nextHeaderKey = nextHeaderKey
	ch.previousChainMessagesCount = ch.nextMessageNumber
	ch.nextMessageNumber = 0
}

func (ch *Chain) advance() (messageKey keys.Message, err error) {
	if ch.masterKey == nil {
		return messageKey, fmt.Errorf("%w: master key is nil", errlist.ErrInvalidValue)
	}

	newMasterKey, messageKey, err := ch.cfg.crypto.AdvanceChain(*ch.masterKey)
	if err != nil {
		return messageKey, fmt.Errorf("%w: advance via crypto: %w", errlist.ErrCrypto, err)
	}

	ch.masterKey = &newMasterKey
	ch.nextMessageNumber++

	return messageKey, err
}
