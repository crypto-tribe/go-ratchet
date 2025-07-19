package rootchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crypto-tribe/go-ratchet/keys"
)

type testCrypto struct{}

func (testCrypto) AdvanceChain(
	_ keys.Root,
	_ keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	return keys.Root{}, keys.Master{}, keys.Header{}, nil
}

var newConfigTests = []struct {
	name           string
	options        []Option
	errCategories  []error
	expectedCrypto Crypto
}{
	{
		"default",
		nil,
		nil,
		defaultCrypto{},
	},
	{
		"crypto option success",
		[]Option{
			WithCrypto(testCrypto{}),
		},
		nil,
		testCrypto{},
	},
	{
		"nil crypto",
		[]Option{
			WithCrypto(nil),
		},
		[]error{
			ErrApplyOptions,
			ErrCryptoIsNil,
		},
		nil,
	},
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	for _, test := range newConfigTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := newConfig(test.options...)
			if err != nil && len(test.errCategories) == 0 {
				t.Fatalf("newConfig() expected no error but got %v", err)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf("newConfig() expected error %v but got %v", errCategory, err)
				}
			}

			if err != nil {
				return
			}

			if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(test.expectedCrypto) {
				t.Fatal("WithCrypto() option did not set passed crypto")
			}
		})
	}
}
