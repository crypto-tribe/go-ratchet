package rootchain

import (
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
)

func TestDefaultCryptoAdvanceChain(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	tests := []struct {
		name      string
		rootKey   keys.Root
		sharedKey keys.Shared
	}{
		{"zero keys", keys.Root{}, keys.Shared{}},
		{"full shared key and zero root key", keys.Root{}, keys.Shared{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8}}},
		{"full root key and zero shared key", keys.Root{Bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}}, keys.Shared{}},
		{
			"full keys",
			keys.Root{Bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}},
			keys.Shared{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rootKey, messageMasterKey, headerKey, err := crypto.AdvanceChain(test.rootKey, test.sharedKey)
			if err != nil {
				t.Fatalf("AdvanceChain(%+v, %+v): expected no error but got %v", test.rootKey, test.sharedKey, err)
			}

			if reflect.DeepEqual(rootKey, test.rootKey) {
				t.Fatalf("AdvanceChain(%+v, %+v): expected different root key", test.rootKey, test.sharedKey)
			}

			if len(rootKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v, %+v): returned empty root key", test.rootKey, test.sharedKey)
			}

			if len(messageMasterKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v, %+v): returned empty message master key", test.rootKey, test.sharedKey)
			}

			if len(headerKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v, %+v): returned empty header key", test.rootKey, test.sharedKey)
			}
		})
	}
}
