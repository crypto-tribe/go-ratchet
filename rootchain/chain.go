package rootchain

import (
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

type Chain struct {
	rootKey keys.Root
	cfg     config
}

func New(rootKey keys.Root, options ...Option) (Chain, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	return Chain{rootKey: rootKey, cfg: cfg}, nil
}

func (ch *Chain) Advance(sharedKey keys.Shared) (keys.MessageMaster, keys.Header, error) {
	newRootKey, messageMasterKey, nextHeaderKey, err := ch.cfg.crypto.AdvanceChain(ch.rootKey, sharedKey)
	if err != nil {
		return keys.MessageMaster{}, keys.Header{}, fmt.Errorf("%w: advance: %w", errlist.ErrCrypto, err)
	}

	ch.rootKey = newRootKey

	return messageMasterKey, nextHeaderKey, nil
}

func (ch Chain) Clone() Chain {
	ch.rootKey = ch.rootKey.Clone()
	return ch
}
