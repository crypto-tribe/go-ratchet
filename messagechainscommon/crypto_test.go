package messagechainscommon

import (
	"testing"

	"github.com/platform-inf/go-ratchet/keys"
)

func TestDeriveMessageCipherKeyAndNonce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		messageKey keys.Message
		errString  string
	}{
		{},
		{keys.Message{Bytes: []byte{1, 2, 3}}, ""},
	}

	for _, test := range tests {
		_, _, err := DeriveMessageCipherKeyAndNonce(test.messageKey)
		if (err == nil && test.errString != "") || (err != nil && err.Error() != test.errString) {
			t.Fatalf("DeriveMessageCipherKeyAndNonce(%+v): expected error %q but got %v", test.messageKey, test.errString, err)
		}
	}
}
