package rootchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestNewChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		rootKey       keys.Root
		options       []Option
		errCategories []error
		errString     string
	}{
		{"zero root key and no options", keys.Root{}, nil, nil, ""},
		{
			"full root key and crypto option",
			keys.Root{Bytes: []byte{1, 2, 3}},
			[]Option{WithCrypto(testCrypto{})},
			nil,
			"",
		},
		{
			"zero key and crypto option error",
			keys.Root{},
			[]Option{WithCrypto(nil)},
			[]error{errlist.ErrInvalidValue, errlist.ErrOption},
			"new config: option: invalid value: crypto is nil",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(test.rootKey, test.options...)
			if (err == nil && test.errString != "") || (err != nil && err.Error() != test.errString) {
				t.Fatalf("New(%+v, %+v): expected error %q but got %+v", test.rootKey, test.options, test.errString, err)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf("New(%+v, %+v): expected error category %v but got %v", test.rootKey, test.options, errCategory, err)
				}
			}

			if !reflect.DeepEqual(chain.rootKey, test.rootKey) {
				t.Fatalf("New(%+v, %+v): invalid root key: %v != %v", test.rootKey, test.options, chain.rootKey, test.rootKey)
			}
		})
	}
}

func TestChainAdvance(t *testing.T) {
	t.Parallel()

	chain, err := New(keys.Root{Bytes: []byte{1, 2, 3, 4, 5}})
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	tests := []struct {
		name      string
		sharedKey keys.Shared
	}{
		{"zero key", keys.Shared{}},
		{"full key", keys.Shared{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			messageMasterKey, nextHeaderKey, err := chain.Advance(test.sharedKey)
			if err != nil {
				t.Fatalf("%+v.Advance(%+v): expected no error but got %v", chain, test.sharedKey, err)
			}

			if len(messageMasterKey.Bytes) == 0 {
				t.Fatalf("%+v.Advance(%+v): returned empty message master key %v", chain, test.sharedKey, messageMasterKey)
			}

			if len(nextHeaderKey.Bytes) == 0 {
				t.Fatalf("%+v.Advance(%+v): returned empty next header key %v", chain, test.sharedKey, nextHeaderKey)
			}
		})
	}
}

func TestChainClone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rootKey keys.Root
	}{
		{"zero key", keys.Root{}},
		{"full key", keys.Root{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(test.rootKey)
			if err != nil {
				t.Fatalf("New(): expected no error but got %v", err)
			}

			clone := chain.Clone()
			if !reflect.DeepEqual(clone.rootKey, chain.rootKey) {
				t.Fatalf("%+v.Clone(): clone contains different root key: %+v != %+v", chain, clone.rootKey, chain.rootKey)
			}

			if len(chain.rootKey.Bytes) > 0 && &chain.rootKey.Bytes[0] == &clone.rootKey.Bytes[0] {
				t.Fatalf("%+v.Clone(): clone contains same root key memory pointer %v", chain, clone.rootKey)
			}
		})
	}
}
