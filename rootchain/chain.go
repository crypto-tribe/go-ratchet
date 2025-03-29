package rootchain

import (
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

type Chain struct {
	key keys.Root
	cfg config
}

func New(key keys.Root, options ...Option) (Chain, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return Chain{}, fmt.Errorf("new config: %w", err)
	}

	return Chain{key: key, cfg: cfg}, nil
}

func (ch *Chain) Advance(sharedKey keys.Shared) (keys.MessageMaster, keys.Header, error) {
	newRootKey, messageMasterKey, nextHeaderKey, err := ch.cfg.crypto.AdvanceChain(ch.key, sharedKey)
	if err != nil {
		return keys.MessageMaster{}, keys.Header{}, fmt.Errorf("%w: advance: %w", errlist.ErrCrypto, err)
	}

	ch.key = newRootKey

	return messageMasterKey, nextHeaderKey, nil
}

func (ch Chain) Clone() Chain {
	ch.key = ch.key.Clone()
	return ch
}
