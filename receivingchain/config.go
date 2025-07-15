package receivingchain

import (
	"errors"

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
		return config{}, errors.Join(ErrApplyOptions, err)
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

// Option is the way to modify config default values.
type Option func(cfg *config) error

// WithCrypto sets passed crypto to the config.
func WithCrypto(crypto Crypto) Option {
	return func(cfg *config) error {
		if check.IsNil(crypto) {
			return ErrCryptoIsNil
		}

		cfg.crypto = crypto

		return nil
	}
}

// WithSkippedKeysStorage sets passed storage to the config.
func WithSkippedKeysStorage(storage SkippedKeysStorage) Option {
	return func(cfg *config) (err error) {
		if check.IsNil(storage) {
			return ErrSkippedKeysStorageIsNil
		}

		cfg.skippedKeysStorage = storage

		return err
	}
}
