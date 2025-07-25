package sendingchain

import (
	"errors"

	"github.com/platform-source/aegis/header"
	"github.com/platform-source/aegis/keys"
	"github.com/platform-source/tools/convert"
	"github.com/platform-source/tools/slices"
)

// Chain is the ratchet sending chain.
//
// Please note that this structure may corrupt its state in case of errors.
// Therefore, clone the data at the top level and replace the current data
// with it if there are no errors.
type Chain struct {
	masterKey                  *keys.Master
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	cfg                        config
}

// New creates a new sending chain.
func New(
	masterKey *keys.Master,
	headerKey *keys.Header,
	nextHeaderKey keys.Header,
	nextMessageNumber uint64,
	previousChainMessagesCount uint64,
	options ...Option,
) (Chain, error) {
	chain := Chain{
		masterKey:                  masterKey,
		headerKey:                  headerKey,
		nextHeaderKey:              nextHeaderKey,
		nextMessageNumber:          nextMessageNumber,
		previousChainMessagesCount: previousChainMessagesCount,
	}

	var err error

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return Chain{}, errors.Join(ErrNewConfig, err)
	}

	return chain, nil
}

// Clone clones sending chain.
func (ch Chain) Clone() Chain {
	ch.masterKey = ch.masterKey.ClonePtr()
	ch.headerKey = ch.headerKey.ClonePtr()
	ch.nextHeaderKey = ch.nextHeaderKey.Clone()

	return ch
}

// Encrypt encrypts passed header and data and authenticates with passed auth.
func (ch *Chain) Encrypt(
	head header.Header,
	data []byte,
	auth []byte,
) (encryptedHeader []byte, encryptedData []byte, err error) {
	if ch.headerKey == nil {
		return nil, nil, ErrHeaderKeyIsNil
	}

	encryptedHeader, err = ch.cfg.crypto.EncryptHeader(*ch.headerKey, head)
	if err != nil {
		return nil, nil, errors.Join(ErrEncryptHeader, err)
	}

	messageKey, err := ch.advance()
	if err != nil {
		return nil, nil, errors.Join(ErrAdvanceChain, err)
	}

	auth = slices.ConcatBytes(encryptedHeader, auth)

	encryptedData, err = ch.cfg.crypto.EncryptMessage(messageKey, data, auth)
	if err != nil {
		return nil, nil, errors.Join(ErrEncryptMessage, err)
	}

	return encryptedHeader, encryptedData, nil
}

// PrepareHeader prepares a new header to send.
func (ch *Chain) PrepareHeader(publicKey keys.Public) header.Header {
	head := header.Header{
		PublicKey:                         publicKey,
		MessageNumber:                     ch.nextMessageNumber,
		PreviousSendingChainMessagesCount: ch.previousChainMessagesCount,
	}

	return head
}

// Upgrade upgrades sending chain with new starting values.
func (ch *Chain) Upgrade(masterKey keys.Master, nextHeaderKey keys.Header) {
	ch.masterKey = &masterKey
	ch.headerKey = convert.ToPtr(ch.nextHeaderKey.Clone())
	ch.nextHeaderKey = nextHeaderKey
	ch.previousChainMessagesCount = ch.nextMessageNumber
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
