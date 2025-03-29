package header

import (
	"errors"
	"slices"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestHeaderSuccessEncodeAndDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		header Header
		bytes  []byte
	}{
		{
			Header{
				PublicKey:                         keys.Public{Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05}},
				PreviousSendingChainMessagesCount: 123,
				MessageNumber:                     321,
			},
			[]byte{
				0x41, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x02, 0x03, 0x04, 0x05,
			},
		},
		{
			Header{},
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, test := range tests {
		bytes := test.header.Encode()
		if !slices.Equal(bytes, test.bytes) {
			t.Fatalf("%+v.Encode(): expected %v but got %v", test.header, test.bytes, bytes)
		}

		header, err := Decode(bytes)
		if err != nil {
			t.Fatalf("Decode(%v): expected no error but got %v", bytes, err)
		}

		if !slices.Equal(header.PublicKey.Bytes, test.header.PublicKey.Bytes) {
			t.Fatalf("Decode(%v): invalid public key: %v != %v", bytes, header.PublicKey.Bytes, test.header.PublicKey.Bytes)
		}

		if header.PreviousSendingChainMessagesCount != test.header.PreviousSendingChainMessagesCount {
			t.Fatalf(
				"Decode(%v): invalid previous sending chain message count: %v != %v",
				bytes,
				header.PreviousSendingChainMessagesCount,
				test.header.PreviousSendingChainMessagesCount,
			)
		}

		if header.MessageNumber != test.header.MessageNumber {
			t.Fatalf("Decode(%v): invalid message number: %v != %v", bytes, header.MessageNumber, test.header.MessageNumber)
		}
	}
}

func TestHeaderDecodeError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		bytes         []byte
		errorCategory error
		errorString   string
	}{
		{
			[]byte{
				0x12, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x0F,
				0x55, 0x00, 0x00, 0x00, 0x77, 0x00, 0x0B,
			},
			errlist.ErrInvalidValue,
			"invalid value: not enough bytes",
		},
		{nil, errlist.ErrInvalidValue, "invalid value: not enough bytes"},
	}

	for _, test := range tests {
		if _, err := Decode(test.bytes); !errors.Is(err, test.errorCategory) || err.Error() != test.errorString {
			t.Fatalf("Decode(%v) expected error %q but got %v", test.bytes, test.errorString, err)
		}
	}
}
