package rootchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-utils/check"
)

type config struct {
	crypto Crypto
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

// Option is the way to modify config default values.
type Option func(cfg *config) error

// WithCrypto is an option to set specific crypto to the config.
func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if check.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}
