package sendingchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-utils"
)

type config struct {
	crypto Crypto
}

func newConfig(options ...Option) (cfg config, err error) {
	cfg = config{
		crypto: newDefaultCrypto(),
	}

	err = cfg.applyOptions(options...)
	if err != nil {
		return cfg, fmt.Errorf("%w: %w", errlist.ErrOption, err)
	}

	return cfg, err
}

func (cfg *config) applyOptions(options ...Option) (err error) {
	for _, option := range options {
		err = option(cfg)
		if err != nil {
			return err
		}
	}

	return err
}

type Option func(cfg *config) error

func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) (err error) {
		if utils.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return err
	}
}
