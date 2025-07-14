package chainscommon

import (
	"testing"

	"github.com/lyreware/go-ratchet/keys"
)

var deriveMessageCipherKeyAndNonceTests = []struct {
	name       string
	messageKey keys.Message
}{
	{
		"zero key",
		keys.Message{},
	},
	{
		"empty key",
		keys.Message{
			Bytes: []byte{},
		},
	},
	{
		"non-empty key",
		keys.Message{
			Bytes: []byte{1, 2, 3},
		},
	},
}

func TestDeriveMessageCipherKeyAndNonce(t *testing.T) {
	t.Parallel()

	for _, test := range deriveMessageCipherKeyAndNonceTests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cipherKey, cipherNonce, err := DeriveMessageCipherKeyAndNonce(test.messageKey)
			if err != nil {
				t.Fatalf(
					"DeriveMessageCipherKeyAndNonce(%+v): expected no error but got %v",
					test.messageKey,
					err,
				)
			}

			if len(cipherKey) == 0 {
				t.Fatalf(
					"DeriveMessageCipherKeyAndNonce(%+v): returned empty cipher key",
					test.messageKey,
				)
			}

			if len(cipherNonce) == 0 {
				t.Fatalf(
					"DeriveMessageCipherKeyAndNonce(%+v): returned empty cipher nonce",
					test.messageKey,
				)
			}
		})
	}
}
