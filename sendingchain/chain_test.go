package sendingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestNewChain(t *testing.T) {
	t.Parallel()

	type args struct {
		masterKey                  *keys.MessageMaster
		headerKey                  *keys.Header
		nextHeaderKey              keys.Header
		nextMessageNumber          uint64
		previousChainMessagesCount uint64
		options                    []Option
	}

	tests := []struct {
		name          string
		args          args
		errCategories []error
		errString     string
	}{
		{"zero args and no options", args{}, nil, ""},
		{
			"full args and crypto option",
			args{
				&keys.MessageMaster{Bytes: []byte{1, 2, 3}},
				&keys.Header{Bytes: []byte{4, 5, 6}},
				keys.Header{Bytes: []byte{7, 8, 9}},
				12,
				201,
				[]Option{WithCrypto(testCrypto{})},
			},
			nil,
			"",
		},
		{
			"zero args and crypto option error",
			args{options: []Option{WithCrypto(nil)}},
			[]error{errlist.ErrInvalidValue, errlist.ErrOption},
			"new config: option: invalid value: crypto is nil",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(
				test.args.masterKey,
				test.args.headerKey,
				test.args.nextHeaderKey,
				test.args.nextMessageNumber,
				test.args.previousChainMessagesCount,
				test.args.options...,
			)
			if (err == nil && test.errString != "") || (err != nil && err.Error() != test.errString) {
				t.Fatalf("New(%+v): expected error %q but got %+v", test.args, test.errString, err)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf("New(%+v): expected error category %v but got %v", test.args, errCategory, err)
				}
			}

			if !reflect.DeepEqual(chain.masterKey, test.args.masterKey) {
				t.Fatalf("New(%+v): invalid master key: %v != %v", test.args, test.args.masterKey, chain.masterKey)
			}

			if !reflect.DeepEqual(chain.headerKey, test.args.headerKey) {
				t.Fatalf("New(%+v): invalid header key: %v != %v", test.args, test.args.headerKey, chain.headerKey)
			}

			if !reflect.DeepEqual(chain.nextHeaderKey, test.args.nextHeaderKey) {
				t.Fatalf("New(%+v): invalid next header key: %v != %v", test.args, test.args.nextHeaderKey, chain.nextHeaderKey)
			}

			if chain.nextMessageNumber != test.args.nextMessageNumber {
				t.Fatalf(
					"New(%+v): invalid message number: %v != %v", test.args, test.args.nextMessageNumber, chain.nextMessageNumber)
			}

			if chain.previousChainMessagesCount != test.args.previousChainMessagesCount {
				t.Fatalf(
					"New(%+v): invalid message number: %v != %v",
					test.args,
					test.args.previousChainMessagesCount,
					chain.previousChainMessagesCount,
				)
			}
		})
	}
}

func TestChainClone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		masterKey                  *keys.MessageMaster
		headerKey                  *keys.Header
		nextHeaderKey              keys.Header
		nextMessageNumber          uint64
		previousChainMessagesCount uint64
		options                    []Option
	}{
		{name: "zero args"},
		{
			"full args",
			&keys.MessageMaster{Bytes: []byte{1, 2, 3}},
			&keys.Header{Bytes: []byte{4, 5, 6}},
			keys.Header{Bytes: []byte{7, 8, 9}},
			12,
			201,
			[]Option{WithCrypto(testCrypto{})},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(
				test.masterKey,
				test.headerKey,
				test.nextHeaderKey,
				test.nextMessageNumber,
				test.previousChainMessagesCount,
				test.options...,
			)
			if err != nil {
				t.Fatalf("New(): expected no error but got %v", err)
			}

			clone := chain.Clone()

			if !reflect.DeepEqual(clone.masterKey, chain.masterKey) {
				t.Fatalf("%+v.Clone(): clone contains different master key: %+v", chain, clone.masterKey)
			}

			if chain.masterKey != nil &&
				len(chain.masterKey.Bytes) > 0 &&
				&chain.masterKey.Bytes[0] == &clone.masterKey.Bytes[0] {
				t.Fatalf("%+v.Clone(): clone contains same master key memory pointer %v", chain, clone.masterKey)
			}

			if !reflect.DeepEqual(clone.headerKey, chain.headerKey) {
				t.Fatalf("%+v.Clone(): clone contains different header key: %+v", chain, clone.headerKey)
			}

			if chain.headerKey != nil &&
				len(chain.headerKey.Bytes) > 0 &&
				&chain.headerKey.Bytes[0] == &clone.headerKey.Bytes[0] {
				t.Fatalf("%+v.Clone(): clone contains same header key memory pointer %v", chain, clone.headerKey)
			}

			if !reflect.DeepEqual(clone.nextHeaderKey, chain.nextHeaderKey) {
				t.Fatalf("%+v.Clone(): clone contains different next header key: %+v", chain, clone.nextHeaderKey)
			}

			if len(chain.nextHeaderKey.Bytes) > 0 && &chain.nextHeaderKey.Bytes[0] == &clone.nextHeaderKey.Bytes[0] {
				t.Fatalf("%+v.Clone(): clone contains same next header key key memory pointer %v", chain, clone.nextHeaderKey)
			}
		})
	}
}

func TestChainPrepareHeader(t *testing.T) {
	t.Parallel()

	type args struct {
		publicKey                  keys.Public
		nextMessageNumber          uint64
		previousChainMessagesCount uint64
	}

	tests := []struct {
		name   string
		args   args
		header header.Header
	}{
		{"zero args and header", args{}, header.Header{}},
		{
			"full args and header",
			args{
				keys.Public{Bytes: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
				123,
				456,
			},
			header.Header{
				PublicKey:                         keys.Public{Bytes: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
				MessageNumber:                     123,
				PreviousSendingChainMessagesCount: 456,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(nil, nil, keys.Header{}, test.args.nextMessageNumber, test.args.previousChainMessagesCount)
			if err != nil {
				t.Fatalf("New(): expected no error but got %v", err)
			}

			header := chain.PrepareHeader(test.args.publicKey)
			if !reflect.DeepEqual(header, test.header) {
				t.Fatalf("%+v.PrepareHeader(%+v): expected %+v but got %+v", chain, test.args.publicKey, test.header, header)
			}
		})
	}
}

func TestChainUpgrade(t *testing.T) {
	t.Parallel()

	oldNextHeaderKey := keys.Header{Bytes: []byte{1, 2, 3}}
	var oldNextMessageNumber uint64 = 222

	chain, err := New(nil, nil, oldNextHeaderKey, oldNextMessageNumber, 111)
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	masterKey := keys.MessageMaster{Bytes: []byte{11, 22, 33}}
	nextHeaderKey := keys.Header{Bytes: []byte{44, 55, 66, 77}}

	chain.Upgrade(masterKey, nextHeaderKey)

	if !reflect.DeepEqual(*chain.masterKey, masterKey) {
		t.Fatalf("Upgrade(%+v, %+v): set different master key %+v", masterKey, nextHeaderKey, *chain.masterKey)
	}

	if !reflect.DeepEqual(*chain.headerKey, oldNextHeaderKey) {
		t.Fatalf(
			"Upgrade(%+v, %+v): header key %+v is not old next header key %+v",
			masterKey,
			nextHeaderKey,
			*chain.headerKey,
			oldNextHeaderKey,
		)
	}

	if !reflect.DeepEqual(chain.nextHeaderKey, nextHeaderKey) {
		t.Fatalf("Upgrade(%+v, %+v): set different next header key %+v", masterKey, nextHeaderKey, chain.nextHeaderKey)
	}

	if chain.nextMessageNumber != 0 {
		t.Fatalf("Upgrade(%+v, %+v): expected message number 0 but got %d", masterKey, nextHeaderKey, chain.nextMessageNumber)
	}

	if chain.previousChainMessagesCount != oldNextMessageNumber {
		t.Fatalf(
			"Upgrade(%+v, %+v): expected previous count %d but got %d",
			masterKey,
			nextHeaderKey,
			chain.previousChainMessagesCount,
			oldNextMessageNumber,
		)
	}
}
