package rootchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crypto-tribe/go-ratchet/keys"
)

var newChainTests = []struct {
	name          string
	rootKey       keys.Root
	options       []Option
	errCategories []error
}{
	{
		"zero root key and no options",
		keys.Root{},
		nil,
		nil,
	},
	{
		"non-empty root key and crypto option",
		keys.Root{
			Bytes: []byte{1, 2, 3},
		},
		[]Option{
			WithCrypto(testCrypto{}),
		},
		nil,
	},
	{
		"zero key and nil crypto",
		keys.Root{},
		[]Option{
			WithCrypto(nil),
		},
		[]error{
			ErrNewConfig,
			ErrApplyOptions,
			ErrCryptoIsNil,
		},
	},
}

func TestNewChain(t *testing.T) {
	t.Parallel()

	for _, test := range newChainTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(test.rootKey, test.options...)
			if err != nil && len(test.errCategories) == 0 {
				t.Fatalf(
					"New(%+v, %+v): expected no error but got %v",
					test.rootKey,
					test.options,
					err,
				)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf(
						"New(%+v, %+v): expected error category %v but got %v",
						test.rootKey,
						test.options,
						errCategory,
						err,
					)
				}
			}
		})
	}
}

var chainAdvanceTests = []struct {
	name      string
	sharedKey keys.Shared
}{
	{
		"zero key",
		keys.Shared{},
	},
	{
		"non-empty key",
		keys.Shared{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestChainAdvance(t *testing.T) {
	t.Parallel()

	chain, err := New(
		keys.Root{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	)
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	for _, test := range chainAdvanceTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			masterKey, nextHeaderKey, err := chain.Advance(test.sharedKey)
			if err != nil {
				t.Fatalf(
					"%+v.Advance(%+v): expected no error but got %v",
					chain,
					test.sharedKey,
					err,
				)
			}

			if len(masterKey.Bytes) == 0 {
				t.Fatalf(
					"%+v.Advance(%+v): returned empty message master key %v",
					chain,
					test.sharedKey,
					masterKey,
				)
			}

			if len(nextHeaderKey.Bytes) == 0 {
				t.Fatalf(
					"%+v.Advance(%+v): returned empty next header key %v",
					chain,
					test.sharedKey,
					nextHeaderKey,
				)
			}
		})
	}
}

var chainCloneTests = []struct {
	name    string
	rootKey keys.Root
}{
	{
		"zero key",
		keys.Root{},
	},
	{
		"non-empty key",
		keys.Root{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestChainClone(t *testing.T) {
	t.Parallel()

	for _, test := range chainCloneTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(test.rootKey)
			if err != nil {
				t.Fatalf("New(): expected no error but got %v", err)
			}

			clone := chain.Clone()
			if !reflect.DeepEqual(clone.rootKey, chain.rootKey) {
				t.Fatalf(
					"%+v.Clone(): clone contains different root key: %+v != %+v",
					chain,
					clone.rootKey,
					chain.rootKey,
				)
			}

			if len(chain.rootKey.Bytes) > 0 && &chain.rootKey.Bytes[0] == &clone.rootKey.Bytes[0] {
				t.Fatalf(
					"%+v.Clone(): clone contains same root key memory pointer %v",
					chain,
					clone.rootKey,
				)
			}
		})
	}
}
