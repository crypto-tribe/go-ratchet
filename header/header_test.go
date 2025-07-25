package header

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/platform-source/aegis/keys"
)

var encodeAndDecodeTests = []struct {
	name   string
	header Header
	bytes  []byte
}{
	{
		"zero header",
		Header{},
		[]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
		"non-empty header",
		Header{
			PublicKey: keys.Public{
				Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			},
			PreviousSendingChainMessagesCount: 123,
			MessageNumber:                     321,
		},
		[]byte{
			0x41, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x02, 0x03, 0x04, 0x05,
		},
	},
}

func TestEncodeAndDecode(t *testing.T) {
	t.Parallel()

	for _, test := range encodeAndDecodeTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bytes := test.header.Encode()
			if !slices.Equal(bytes, test.bytes) {
				t.Fatalf("%+v.Encode(): expected %v but got %v", test.header, test.bytes, bytes)
			}

			header, err := Decode(bytes)
			if err != nil {
				t.Fatalf("Decode(%v): expected no error but got %v", bytes, err)
			}

			if !reflect.DeepEqual(header, test.header) {
				t.Fatalf("Decode(%v): expected %+v but got %+v", bytes, test.header, header)
			}
		})
	}
}

var decodeTests = []struct {
	name          string
	bytes         []byte
	errCategories []error
}{
	{
		"not enough bytes",
		[]byte{
			0x12, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x0F,
			0x55, 0x00, 0x00, 0x00, 0x77, 0x00, 0x0B,
		},
		[]error{
			ErrNotEnoughBytes,
		},
	},
	{
		"nil bytes slice",
		nil,
		[]error{
			ErrNotEnoughBytes,
		},
	},
}

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, test := range decodeTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := Decode(test.bytes)

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf(
						"Decode(%v) expected error %q but got %v",
						test.bytes,
						errCategory,
						err,
					)
				}
			}
		})
	}
}
