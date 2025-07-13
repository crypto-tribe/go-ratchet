package receivingchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-utils/check"
)

type config struct {
	crypto             Crypto
	skippedKeysStorage SkippedKeysStorage
}

func newConfig(options ...Option) (config, error) {
	cfg := config{
		crypto:             newDefaultCrypto(),
		skippedKeysStorage: newDefaultSkippedKeysStorage(),
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

func (cfg config) clone() config {
	cfg.skippedKeysStorage = cfg.skippedKeysStorage.Clone()

	return cfg
}

type Option func(cfg *config) error

func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if check.IsNil(crypto) {
			return fmt.Errorf("%w: crypto is nil", errlist.ErrInvalidValue)
		}

		cfg.crypto = crypto

		return nil
	}
}

func WithSkippedKeysStorage(storage SkippedKeysStorage) Option {
	return func(cfg *config) (err error) {
		if check.IsNil(storage) {
			return fmt.Errorf("%w: storage is nil", errlist.ErrInvalidValue)
		}

		cfg.skippedKeysStorage = storage

		return err
	}
}
