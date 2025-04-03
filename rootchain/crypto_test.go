package rootchain

import (
	"slices"
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
)

func TestDefaultCryptoAdvanceChain(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	tests := []struct {
		rootKey   keys.Root
		sharedKey keys.Shared
	}{
		{},
		{sharedKey: keys.Shared{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8}}},
		{rootKey: keys.Root{Bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}}},
		{
			sharedKey: keys.Shared{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8}},
			rootKey:   keys.Root{Bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
	}

	for _, test := range tests {
		rootKey, messageMasterKey, headerKey, err := crypto.AdvanceChain(test.rootKey, test.sharedKey)
		if err != nil {
			t.Fatalf("AdvanceChain(%+v, %+v): expected no error but got %v", test.rootKey, test.sharedKey, err)
		}

		if len(rootKey.Bytes) == 0 {
			t.Fatalf("AdvanceChain(%+v, %+v): returned empty root key", test.rootKey, test.sharedKey)
		}

		if slices.Equal(rootKey.Bytes, test.rootKey.Bytes) {
			t.Fatalf("AdvanceChain(%+v, %+v): expected different root key", test.rootKey, test.sharedKey)
		}

		if len(messageMasterKey.Bytes) == 0 {
			t.Fatalf("AdvanceChain(%+v, %+v): returned empty message master key", test.rootKey, test.sharedKey)
		}

		if len(headerKey.Bytes) == 0 {
			t.Fatalf("AdvanceChain(%+v, %+v): returned empty header key", test.rootKey, test.sharedKey)
		}
	}
}
