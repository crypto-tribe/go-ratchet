package rootchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/keys"
)

type Chain struct {
	rootKey keys.Root
	cfg     config
}

func New(rootKey keys.Root, options ...Option) (chain Chain, err error) {
	chain = Chain{
		rootKey: rootKey,
	}

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return chain, fmt.Errorf("new config: %w", err)
	}

	return chain, err
}

func (ch *Chain) Advance(sharedKey keys.Shared) (masterKey keys.Master, nextHeaderKey keys.Header, err error) {
	ch.rootKey, masterKey, nextHeaderKey, err = ch.cfg.crypto.AdvanceChain(ch.rootKey, sharedKey)
	if err != nil {
		return masterKey, nextHeaderKey, fmt.Errorf("%w: advance: %w", errlist.ErrCrypto, err)
	}

	return masterKey, nextHeaderKey, err
}

func (ch Chain) Clone() Chain {
	ch.rootKey = ch.rootKey.Clone()
	return ch
}
