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
		errString string
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
		rootKey, _, _, err := crypto.AdvanceChain(test.rootKey, test.sharedKey)
		if (err == nil && test.errString != "") || (err != nil && err.Error() != test.errString) {
			t.Fatalf("AdvanceChain(%+v, %+v): expected error %q but got %v", test.rootKey, test.sharedKey, test.errString, err)
		}

		if slices.Equal(rootKey.Bytes, test.rootKey.Bytes) {
			t.Fatalf("AdvanceChain(%+v, %+v): expected different root key", test.rootKey, test.sharedKey)
		}
	}
}
