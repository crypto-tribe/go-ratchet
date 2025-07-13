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

func New(rootKey keys.Root, options ...Option) (Chain, error) {
	chain := Chain{
		rootKey: rootKey,
	}

	var err error

	chain.cfg, err = newConfig(options...)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	return chain, err
}

func (ch *Chain) Advance(
	sharedKey keys.Shared,
) (masterKey keys.Master, nextHeaderKey keys.Header, err error) {
	ch.rootKey, masterKey, nextHeaderKey, err = ch.cfg.crypto.AdvanceChain(ch.rootKey, sharedKey)
	if err != nil {
		return keys.Master{}, keys.Header{}, fmt.Errorf("%w: advance: %w", errlist.ErrCrypto, err)
	}

	return masterKey, nextHeaderKey, nil
}

func (ch Chain) Clone() Chain {
	ch.rootKey = ch.rootKey.Clone()

	return ch
}
