package messagechainscommon

import (
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
)

func TestDeriveMessageCipherKeyAndNonce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		messageKey keys.Message
	}{
		{"zero key", keys.Message{}},
		{"key with empty bytes slice", keys.Message{Bytes: []byte{}}},
		{"full key", keys.Message{Bytes: []byte{1, 2, 3}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

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
		})
	}
}
