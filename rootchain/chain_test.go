package rootchain

import (
	"errors"
	"slices"
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestNewChain(t *testing.T) {
	t.Parallel()

	tests := []struct{
		rootKey   keys.Root
		options   []Option
		errString string
	}{
		{
			keys.Root{},
			nil,
			"",
		},
		{
			keys.Root{Bytes: []byte{1, 2, 3}},
			[]Option{WithCrypto(testCrypto{})},
			"",
		},
		{
			keys.Root{},
			[]Option{WithCrypto(nil)},
			"config: option: crypto is nil",
		},
	}

	for _, test := range tests {
		chain, err := New(test.rootKey, test.options...)
		if (err == nil && test.errString != "") ||
		(err != nil && err.Error() != test.errString && !errors.Is(err, errlist.ErrOption)) {
			t.Fatalf("New(%+v, %+v): expected no error but got %v", test.rootKey, test.options, err)
		}

		if !reflect.DeepEqual(chain.rootKey, test.rootKey) {
			t.Fatalf("New(%+v, %+v): invalid root key: %v != %v", test.rootKey, test.options, chain.rootKey, test.rootKey)
		}
	}
}

func TestChainAdvance(t *testing.T) {
	chain, err := New(keys.Root{Bytes: []byte{1, 2, 3, 4, 5}})
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	tests := []struct{
		sharedKey keys.Shared
	}{
		{},
		{keys.Shared{Bytes: nil}},
		{keys.Shared{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
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
	}
}

func TestChainClone(t *testing.T) {
	t.Parallel()

	chain, err := New(keys.Root{Bytes: []byte{1, 2, 3, 4, 5}})
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	clone := chain.Clone()

	if !slices.Equal(clone.rootKey.Bytes, chain.rootKey.Bytes) {
		t.Fatalf("%+v.Clone(): clone contains different root key: %+v != %+v", chain, clone.rootKey, chain.rootKey)
	}

	if &chain.rootKey.Bytes[0] == &clone.rootKey.Bytes[0] {
		t.Fatalf("%+v.Clone(): clone contains same root key memory pointer %v", chain, clone.rootKey)
	}
}
