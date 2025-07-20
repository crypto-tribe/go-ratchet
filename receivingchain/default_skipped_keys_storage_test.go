package receivingchain

import (
	"errors"
	"testing"

	"github.com/crypto-tribe/go-ratchet/keys"
)

func TestDefaultSkippedKeysStorageAddClear(t *testing.T) {
	t.Parallel()

	storage := newDefaultSkippedKeysStorage()

	for headerNumber := range defaultSkippedKeysStorageHeaderKeysCountToClear {
		err := storage.Add(
			keys.Header{
				Bytes: make([]byte, headerNumber),
			},
			0,
			keys.Message{},
		)
		if err != nil {
			t.Fatalf("Add(%d): expected no error but got %+v", headerNumber, err)
		}
	}

	if storage.getHeaderKeysCount() != defaultSkippedKeysStorageHeaderKeysCountToClear {
		t.Fatal("Add(): early clear")
	}

	err := storage.Add(
		keys.Header{
			Bytes: []byte{1, 2, 3},
		},
		0,
		keys.Message{},
	)
	if err != nil {
		t.Fatalf("Add(1, 2, 3): expected no error but got %+v", err)
	}

	if storage.getHeaderKeysCount() != 1 {
		t.Fatalf("Add(): expected clear but length is %d", storage.getHeaderKeysCount())
	}
}

func TestDefaultSkippedKeysStorageAddTooManyMessageKeys(t *testing.T) {
	t.Parallel()

	storage := newDefaultSkippedKeysStorage()

	for messageNumber := range defaultSkippedKeysStorageMessageKeysCountLimit {
		err := storage.Add(keys.Header{}, uint64(messageNumber), keys.Message{})
		if err != nil {
			t.Fatalf("Add(%d): expected no error but got %+v", messageNumber, err)
		}
	}

	err := storage.Add(
		keys.Header{},
		defaultSkippedKeysStorageMessageKeysCountLimit,
		keys.Message{},
	)
	if !errors.Is(err, ErrTooManySkippedMessageKeys) {
		t.Fatalf(
			"Add(%d): expected error too many message keys error but got %+v",
			defaultSkippedKeysStorageMessageKeysCountLimit,
			err,
		)
	}
}

func TestDefaultSkippedKeysStorageDelete(t *testing.T) {
	t.Parallel()

	headerKey := keys.Header{Bytes: []byte{1, 2, 3}}
	messageNumber := uint64(456)

	storage := newDefaultSkippedKeysStorage()

	err := storage.Add(headerKey, messageNumber, keys.Message{})
	if err != nil {
		t.Fatalf("Add(): expected no error but got %+v", err)
	}

	if storage.getHeaderKeysCount() != 1 || storage.getMessageKeysCount(headerKey) != 1 {
		t.Fatalf("Add(): expected len 1 but got %d", storage.getHeaderKeysCount())
	}

	err = storage.Delete(keys.Header{}, 100)
	if err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	err = storage.Delete(keys.Header{}, messageNumber)
	if err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	err = storage.Delete(headerKey, 100)
	if err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if storage.getHeaderKeysCount() != 1 || storage.getMessageKeysCount(headerKey) != 1 {
		t.Fatalf(
			"Delete(): expected no delete but len is %d:%d",
			storage.getHeaderKeysCount(),
			storage.getMessageKeysCount(headerKey),
		)
	}

	err = storage.Delete(headerKey, messageNumber)
	if err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if storage.getHeaderKeysCount() != 1 || storage.getMessageKeysCount(headerKey) != 0 {
		t.Fatalf(
			"Delete(): expected delete but len is %d:%d",
			storage.getHeaderKeysCount(),
			storage.getMessageKeysCount(headerKey),
		)
	}
}

func TestDefaultSkippedKeysStorageGetIter(t *testing.T) {
	t.Parallel()

	storage := newDefaultSkippedKeysStorage()
	iters := make(map[byte]map[uint64]byte, defaultSkippedKeysStorageHeaderKeysCountToClear)

	for headerKeyByteInt := range defaultSkippedKeysStorageHeaderKeysCountToClear / 2 {
		headerKeyByte := byte(headerKeyByteInt)
		headerKey := keys.Header{Bytes: []byte{headerKeyByte}}

		for messageNumber := range defaultSkippedKeysStorageMessageKeysCountLimit {
			if messageNumber%2 == 1 {
				continue
			}

			messageKeyByte := byte(messageNumber % (0xFF + 1))
			messageKey := keys.Message{Bytes: []byte{messageKeyByte}}

			err := storage.Add(headerKey, uint64(messageNumber), messageKey)
			if err != nil {
				t.Fatalf(
					"Add(%d, %d): expected no error but got %+v",
					headerKeyByteInt,
					messageNumber,
					err,
				)
			}

			if _, exists := iters[headerKeyByte]; !exists {
				iters[headerKeyByte] = make(
					map[uint64]byte,
					defaultSkippedKeysStorageMessageKeysCountLimit/2,
				)
			}

			iters[headerKeyByte][uint64(messageNumber)] = messageKeyByte
		}
	}

	iter, err := storage.GetIter()
	if err != nil {
		t.Fatalf("GetIter(): expected no error but got %+v", err)
	}

	for headerKey, messageNumberKeys := range iter {
		if len(headerKey.Bytes) != 1 {
			t.Fatalf("GetIter(): expected header key length 1 but got %d", len(headerKey.Bytes))
		}

		messageIters, exists := iters[headerKey.Bytes[0]]
		if !exists {
			t.Fatalf("GetIter(): returned invalid header key byte %d", headerKey.Bytes[0])
		}

		for messageNumber, messageKey := range messageNumberKeys {
			if len(messageKey.Bytes) != 1 {
				t.Fatalf(
					"GetIter(): expected message key length 1 but got %d",
					len(messageKey.Bytes),
				)
			}

			messageKeyByte, exists := messageIters[messageNumber]
			if !exists {
				t.Fatalf("GetIter(): returned invalid message number %d", messageNumber)
			}

			if messageKey.Bytes[0] != messageKeyByte {
				t.Fatalf(
					"GetIter(): got invalid message key byte: %d != %d",
					messageKey.Bytes[0],
					messageKeyByte,
				)
			}

			delete(messageIters, messageNumber)
		}

		if len(messageIters) > 0 {
			t.Fatalf(
				"GetIter(): header key byte is %d, %d iterations remain: %+v",
				headerKey.Bytes[0],
				len(messageIters),
				messageIters,
			)
		}

		delete(iters, headerKey.Bytes[0])
	}

	if len(iters) > 0 {
		t.Fatalf(
			"GetIter(): not enough iterations over header key bytes, %d remain: %+v",
			len(iters),
			iters,
		)
	}
}
