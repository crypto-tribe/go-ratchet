package rootchain

import (
	"errors"

	"github.com/lyreware/go-ratchet/keys"
)

// Chain is the ratchet root chain.
type Chain struct {
	rootKey keys.Root
	cfg     config
}

// New creates a new root chain.
func New(rootKey keys.Root, options ...Option) (Chain, error) {
	chain := Chain{
		rootKey: rootKey,
	}

	var err error

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return Chain{}, errors.Join(ErrNewConfig, err)
	}

	return chain, nil
}

// Advance advances root chain and creates a new master key and next header key.
func (ch *Chain) Advance(sharedKey keys.Shared) (keys.Master, keys.Header, error) {
	var (
		masterKey     keys.Master
		nextHeaderKey keys.Header
		err           error
	)

	ch.rootKey, masterKey, nextHeaderKey, err = ch.cfg.crypto.AdvanceChain(ch.rootKey, sharedKey)
	if err != nil {
		return keys.Master{}, keys.Header{}, errors.Join(ErrCryptoAdvance, err)
	}

	return masterKey, nextHeaderKey, nil
}

// Clone clones a root chain.
func (ch Chain) Clone() Chain {
	ch.rootKey = ch.rootKey.Clone()

	return ch
}
