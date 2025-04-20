package receivingchain

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-utils"
)

func TestDefaultSkippedKeysStorageAdd(t *testing.T) {
	t.Parallel()

	t.Run("test clear", func(t *testing.T) {
		t.Parallel()

		storage := newDefaultSkippedKeysStorage()

		for i := range defaultSkippedKeysStorageHeaderKeysLenToClear {
			var bytes [utils.Uint64Size]byte
			binary.LittleEndian.PutUint64(bytes[:], uint64(i))

			if err := storage.Add(keys.Header{Bytes: bytes[:]}, 0, keys.Message{}); err != nil {
				t.Fatalf("Add(%d): expected no error but got %+v", i, err)
			}
		}

		if len(storage) != defaultSkippedKeysStorageHeaderKeysLenToClear {
			t.Fatal("Add(): early clear")
		}

		if err := storage.Add(keys.Header{Bytes: []byte{1, 2, 3}}, 0, keys.Message{}); err != nil {
			t.Fatalf("Add(1, 2, 3): expected no error but got %+v", err)
		}

		if len(storage) != 1 {
			t.Fatalf("Add(): expected clear but length is %d", len(storage))
		}
	})

	t.Run("test too many message keys", func(t *testing.T) {
		t.Parallel()

		storage := newDefaultSkippedKeysStorage()

		for messageNumber := range defaultSkippedKeysStorageMessageKeysLenLimit {
			if err := storage.Add(keys.Header{}, uint64(messageNumber), keys.Message{}); err != nil {
				t.Fatalf("Add(%d): expected no error but got %+v", messageNumber, err)
			}
		}

		errString := fmt.Sprintf(
			"too many message keys: %d >= %d",
			defaultSkippedKeysStorageMessageKeysLenLimit,
			defaultSkippedKeysStorageMessageKeysLenLimit,
		)

		err := storage.Add(keys.Header{}, defaultSkippedKeysStorageMessageKeysLenLimit, keys.Message{})
		if err == nil || err.Error() != errString {
			t.Fatalf("Add(%d): expected error %q but got %+v", defaultSkippedKeysStorageMessageKeysLenLimit, errString, err)
		}
	})
}

func TestDefaultSkippedKeysStorageDelete(t *testing.T) {
	t.Parallel()

	headerKey := keys.Header{Bytes: []byte{1, 2, 3}}
	messageNumber := uint64(456)

	storage := newDefaultSkippedKeysStorage()
	if err := storage.Add(headerKey, messageNumber, keys.Message{}); err != nil {
		t.Fatalf("Add(): expected no error but got %+v", err)
	}

	storageKey := storage.buildKey(headerKey)

	if len(storage) != 1 || len(storage[storageKey]) != 1 {
		t.Fatalf("Add(): expected len 1 but got %d", len(storage))
	}

	if err := storage.Delete(keys.Header{}, 100); err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if err := storage.Delete(keys.Header{}, messageNumber); err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if err := storage.Delete(headerKey, 100); err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if len(storage) != 1 || len(storage[storageKey]) != 1 {
		t.Fatalf("Delete(): expected no delete but len is %d:%d", len(storage), len(storage[storageKey]))
	}

	if err := storage.Delete(headerKey, messageNumber); err != nil {
		t.Fatalf("Delete(): expected no error but got %+v", err)
	}

	if len(storage) != 1 || len(storage[storageKey]) != 0 {
		t.Fatalf("Delete(): expected delete but len is %d:%d", len(storage), len(storage[storageKey]))
	}
}
