package receivingchain

import (
	"github.com/crypto-tribe/go-ratchet/keys"
)

const (
	defaultSkippedKeysStorageHeaderKeysCountToClear = 4
	defaultSkippedKeysStorageMessageKeysCountLimit  = 1024
)

type defaultSkippedKeysStorage struct {
	mapping map[string]map[uint64]keys.Message
}

func newDefaultSkippedKeysStorage() defaultSkippedKeysStorage {
	storage := defaultSkippedKeysStorage{
		mapping: make(map[string]map[uint64]keys.Message),
	}

	return storage
}

func (st defaultSkippedKeysStorage) Add(
	headerKey keys.Header,
	messageNumber uint64,
	messageKey keys.Message,
) error {
	if st.getHeaderKeysCount() >= defaultSkippedKeysStorageHeaderKeysCountToClear {
		st.clear()
	}

	if st.getMessageKeysCount(headerKey) >= defaultSkippedKeysStorageMessageKeysCountLimit {
		return ErrTooManySkippedMessageKeys
	}

	st.addUnsafe(headerKey, messageNumber, messageKey)

	return nil
}

func (st defaultSkippedKeysStorage) Clone() SkippedKeysStorage {
	stClone := defaultSkippedKeysStorage{
		mapping: make(map[string]map[uint64]keys.Message, st.getHeaderKeysCount()),
	}

	for serHeaderKey, messageNumberKeys := range st.mapping {
		stClone.mapping[serHeaderKey] = make(map[uint64]keys.Message, len(messageNumberKeys))

		for messageNumber, messageKey := range messageNumberKeys {
			stClone.mapping[serHeaderKey][messageNumber] = messageKey.Clone()
		}
	}

	return stClone
}

func (st defaultSkippedKeysStorage) Delete(headerKey keys.Header, messageNumber uint64) error {
	serHeaderKey := st.serializeHeaderKey(headerKey)
	delete(st.mapping[serHeaderKey], messageNumber)

	return nil
}

func (st defaultSkippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	iter := func(yield SkippedKeysYield) {
		for serHeaderKey, messageNumberKeys := range st.mapping {
			headerKey := st.deserializeHeaderKey(serHeaderKey)

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

func (st defaultSkippedKeysStorage) addUnsafe(
	headerKey keys.Header,
	messageNumber uint64,
	messageKey keys.Message,
) {
	serHeaderKey := st.serializeHeaderKey(headerKey)

	if _, ok := st.mapping[serHeaderKey]; !ok {
		st.mapping[serHeaderKey] = make(map[uint64]keys.Message)
	}

	st.mapping[serHeaderKey][messageNumber] = messageKey
}

func (st defaultSkippedKeysStorage) clear() {
	clear(st.mapping)
}

func (defaultSkippedKeysStorage) serializeHeaderKey(key keys.Header) string {
	ser := string(key.Bytes)

	return ser
}

func (defaultSkippedKeysStorage) deserializeHeaderKey(value string) keys.Header {
	headerKey := keys.Header{
		Bytes: []byte(value),
	}

	return headerKey
}

func (st defaultSkippedKeysStorage) getHeaderKeysCount() int {
	count := len(st.mapping)

	return count
}

func (st defaultSkippedKeysStorage) getMessageKeysCount(headerKey keys.Header) int {
	serHeaderKey := st.serializeHeaderKey(headerKey)
	count := len(st.mapping[serHeaderKey])

	return count
}
