package receivingchain

import (
	"github.com/platform-source/aegis/keys"
)

type (
	// SkippedKeysIter passes key values to yield loop iteration.
	SkippedKeysIter func(yield SkippedKeysYield)

	// SkippedKeysYield represents loop iteration.
	SkippedKeysYield func(
		headerKey keys.Header,
		messageNumberKeysIter SkippedMessageNumberKeysIter,
	) bool

	// SkippedMessageNumberKeysIter passes message number and message key to yield loop iteration.
	SkippedMessageNumberKeysIter func(yield SkippedMessageNumberKeysYield)

	// SkippedMessageNumberKeysYield represents loop iteration.
	SkippedMessageNumberKeysYield func(number uint64, key keys.Message) bool

	// SkippedKeysStorage is the storage of skipped keys of the receiving chain.
	SkippedKeysStorage interface {
		// Add must add new skipped keys to storage.
		Add(headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error

		// Clone must deep clone a storage.
		Clone() SkippedKeysStorage

		// Delete must delete skipped keys by header key and message number.
		Delete(headerKey keys.Header, messageNumber uint64) error

		// GetIter must return function, which iterates over all skipped keys.
		GetIter() (SkippedKeysIter, error)
	}
)
