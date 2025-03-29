package receivingchain

import (
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-utils"
)

type config struct {
	crypto             Crypto
	skippedKeysStorage SkippedKeysStorage
}

func newConfig(options []Option) (config, error) {
	cfg := config{
		crypto:             newDefaultCrypto(),
		skippedKeysStorage: newDefaultSkippedKeysStorage(),
	}
	if err := cfg.applyOptions(options); err != nil {
		return config{}, fmt.Errorf("%w: %w", errlist.ErrOption, err)
	}

	return cfg, nil
}

func (cfg *config) applyOptions(options []Option) error {
	for _, option := range options {
		if err := option(cfg); err != nil {
			return err
		}
	}

	return nil
}

func (cfg config) clone() config {
	cfg.skippedKeysStorage = cfg.skippedKeysStorage.Clone()
	return cfg
}

type Option func(cfg *config) error

func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if utils.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}

func WithSkippedKeysStorage(storage SkippedKeysStorage) Option {
	return func(cfg *config) error {
		if utils.IsNil(storage) {
			return fmt.Errorf("%w: storage is nil", errlist.ErrInvalidValue)
		}

		cfg.skippedKeysStorage = storage

		return nil
	}
}
