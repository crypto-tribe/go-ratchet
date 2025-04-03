package messagechainscommon

import (
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
)

func TestDeriveMessageCipherKeyAndNonce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		messageKey keys.Message
	}{
		{},
		{keys.Message{Bytes: []byte{1, 2, 3}}},
	}

	for _, test := range tests {
		cipherKey, cipherNonce, err := DeriveMessageCipherKeyAndNonce(test.messageKey)
		if err != nil {
			t.Fatalf("DeriveMessageCipherKeyAndNonce(%+v): expected no error but got %v", test.messageKey, err)
		}

		if len(cipherKey) == 0 {
			t.Fatalf("DeriveMessageCipherKeyAndNonce(%+v): returned empty cipher key", test.messageKey)
		}

		if len(cipherNonce) == 0 {
			t.Fatalf("DeriveMessageCipherKeyAndNonce(%+v): returned empty cipher nonce", test.messageKey)
		}
	}
}
