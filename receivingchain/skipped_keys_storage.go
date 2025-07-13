package receivingchain

import (
	"fmt"

	"github.com/lyreware/go-ratchet/keys"
)

const (
	defaultSkippedKeysStorageMessageKeysLenLimit  = 1024
	defaultSkippedKeysStorageHeaderKeysLenToClear = 4
)

type (
	// SkippedKeysIter passes key values to yield loop iteration.
	SkippedKeysIter func(yield SkippedKeysYield)

	// SkippedKeysYield represents loop iteration.
	SkippedKeysYield func(headerKey keys.Header, messageNumberKeysIter SkippedMessageNumberKeysIter) bool

	// SkippedMessageNumberKeysIter passes message number and message key to yield loop iteration.
	SkippedMessageNumberKeysIter func(yield SkippedMessageNumberKeysYield)

	// SkippedMessageNumberKeysYield represents loop iteration.
	SkippedMessageNumberKeysYield func(number uint64, key keys.Message) bool
)

// SkippedKeysStorage is the storage of skipped keys of the receiving chain.
type SkippedKeysStorage interface {
	// Add must add new skipped keys to storage.
	Add(headerKey keys.Header, messageNumber uint64, messageKey keys.Message) error

	// Clone must deep clone a storage.
	Clone() SkippedKeysStorage

	// Delete must delete skipped keys by header key and message number.
	Delete(headerKey keys.Header, messageNumber uint64) error

	// GetIter must return function, which iterates over all skipped keys.
	GetIter() (SkippedKeysIter, error)
}

type defaultSkippedKeysStorage map[string]map[uint64]keys.Message

func newDefaultSkippedKeysStorage() defaultSkippedKeysStorage {
	return make(defaultSkippedKeysStorage)
}

func (st defaultSkippedKeysStorage) Add(
	headerKey keys.Header,
	messageNumber uint64,
	messageKey keys.Message,
) error {
	if len(st) >= defaultSkippedKeysStorageHeaderKeysLenToClear {
		clear(st)
	}

	stKey := st.convertToKey(headerKey)
	if len(st[stKey]) >= defaultSkippedKeysStorageMessageKeysLenLimit {
		return fmt.Errorf(
			"too many message keys: %d >= %d",
			len(st[stKey]),
			defaultSkippedKeysStorageMessageKeysLenLimit,
		)
	}

	if _, ok := st[stKey]; !ok {
		st[stKey] = make(map[uint64]keys.Message)
	}

	st[stKey][messageNumber] = messageKey

	return nil
}

func (st defaultSkippedKeysStorage) Clone() SkippedKeysStorage {
	stClone := make(defaultSkippedKeysStorage, len(st))

	for stKey, messageNumberKeys := range st {
		messageNumberKeysClone := make(map[uint64]keys.Message, len(messageNumberKeys))

		for messageNumber, messageKey := range messageNumberKeys {
			messageNumberKeysClone[messageNumber] = messageKey.Clone()
		}

		stClone[stKey] = messageNumberKeysClone
	}

	return stClone
}

func (st defaultSkippedKeysStorage) Delete(headerKey keys.Header, messageNumber uint64) error {
	stKey := st.convertToKey(headerKey)
	delete(st[stKey], messageNumber)

	return nil
}

func (st defaultSkippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	iter := func(yield SkippedKeysYield) {
		for stKey, messageNumberKeys := range st {
			headerKey := st.convertFromKey(stKey)

			messageNumberKeysIter := func(yield SkippedMessageNumberKeysYield) {
				for messageNumber, messageKey := range messageNumberKeys {
					if !yield(messageNumber, messageKey) {
						return
					}
				}
			}

			if !yield(headerKey, messageNumberKeysIter) {
				return
			}
		}
	}

	return iter, nil
}

func (defaultSkippedKeysStorage) convertToKey(headerKey keys.Header) (key string) {
	key = string(headerKey.Bytes)

	return key
}

func (defaultSkippedKeysStorage) convertFromKey(key string) (headerKey keys.Header) {
	headerKey.Bytes = []byte(key)

	return headerKey
}
