package ratchet

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/receivingchain"
	"github.com/lyreware/go-ratchet/rootchain"
	"github.com/lyreware/go-ratchet/sendingchain"
	"github.com/lyreware/go-utils/check"
)

type config struct {
	crypto           Crypto
	receivingOptions []receivingchain.Option
	rootOptions      []rootchain.Option
	sendingOptions   []sendingchain.Option
}

func newConfig(options ...Option) (config, error) {
	cfg := config{
		crypto: newDefaultCrypto(),
	}

	err := cfg.applyOptions(options...)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", errlist.ErrOption, err)
	}

	return cfg, nil
}

func (cfg *config) applyOptions(options ...Option) error {
	for _, option := range options {
		err := option(cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

// Option is a way to modify config default values.
type Option func(cfg *config) error

// WithCrypto sets passed crypto to the config.
func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if check.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}

// WithReceivingChainOptions sets passed options to the receiving chain.
func WithReceivingChainOptions(options ...receivingchain.Option) Option {
	return func(cfg *config) error {
		cfg.receivingOptions = options

		return nil
	}
}

// WithRootChainOptions sets passed options to the root chain.
func WithRootChainOptions(options ...rootchain.Option) Option {
	return func(cfg *config) error {
		cfg.rootOptions = options

		return nil
	}
}

// WithSendingChainOptions sets passed options to the sending chain.
func WithSendingChainOptions(options ...sendingchain.Option) Option {
	return func(cfg *config) error {
		cfg.sendingOptions = options

		return nil
	}
}
