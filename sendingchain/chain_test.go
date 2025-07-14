package sendingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	cipher "golang.org/x/crypto/chacha20poly1305"
)

type newChainTestArgs struct {
	masterKey                  *keys.Master
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	options                    []Option
}

var newChainTests = []struct {
	name          string
	args          newChainTestArgs
	errCategories []error
}{
	{"zero args and no options", newChainTestArgs{}, nil, ""},
	{
		"non-empty args and crypto option",
		newChainTestArgs{
			&keys.Master{
				Bytes: []byte{1, 2, 3},
			},
			&keys.Header{
				Bytes: []byte{4, 5, 6},
			},
			keys.Header{
				Bytes: []byte{7, 8, 9},
			},
			12,
			201,
			[]Option{
				WithCrypto(testCrypto{}),
			},
		},
		nil,
	},
	{
		"zero args and crypto option error",
		newChainTestArgs{
			options: []Option{
				WithCrypto(nil),
			},
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
		test := test

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
			if err != nil && len(test.errCategories) == 0 {
				t.Fatalf("New(%+v): expected no error but got %v", test.args, err)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf(
						"New(%+v): expected error category %v but got %v",
						test.args,
						errCategory,
						err,
					)
				}
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(chain.masterKey, test.args.masterKey) {
				t.Fatalf(
					"New(%+v): invalid master key: %v != %v",
					test.args,
					test.args.masterKey,
					chain.masterKey,
				)
			}

			if !reflect.DeepEqual(chain.headerKey, test.args.headerKey) {
				t.Fatalf(
					"New(%+v): invalid header key: %v != %v",
					test.args,
					test.args.headerKey,
					chain.headerKey,
				)
			}

			if !reflect.DeepEqual(chain.nextHeaderKey, test.args.nextHeaderKey) {
				t.Fatalf(
					"New(%+v): invalid next header key: %v != %v",
					test.args,
					test.args.nextHeaderKey,
					chain.nextHeaderKey,
				)
			}

			if chain.nextMessageNumber != test.args.nextMessageNumber {
				t.Fatalf(
					"New(%+v): invalid message number: %v != %v",
					test.args,
					test.args.nextMessageNumber,
					chain.nextMessageNumber,
				)
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

var chainCloneTests = []struct {
	name                       string
	masterKey                  *keys.Master
	headerKey                  *keys.Header
	nextHeaderKey              keys.Header
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
	options                    []Option
}{
	{name: "zero args"},
	{
		"non-empty args",
		&keys.Master{
			Bytes: []byte{1, 2, 3},
		},
		&keys.Header{
			Bytes: []byte{4, 5, 6},
		},
		keys.Header{
			Bytes: []byte{7, 8, 9},
		},
		12,
		201,
		[]Option{
			WithCrypto(testCrypto{}),
		},
	},
}

func TestChainClone(t *testing.T) {
	t.Parallel()

	for _, test := range chainCloneTests {
		test := test

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
				t.Fatalf(
					"%+v.Clone(): clone contains different master key: %+v",
					chain,
					clone.masterKey,
				)
			}

			if chain.masterKey != nil &&
				len(chain.masterKey.Bytes) > 0 &&
				&chain.masterKey.Bytes[0] == &clone.masterKey.Bytes[0] {
				t.Fatalf(
					"%+v.Clone(): clone contains same master key memory pointer %v",
					chain,
					clone.masterKey,
				)
			}

			if !reflect.DeepEqual(clone.headerKey, chain.headerKey) {
				t.Fatalf(
					"%+v.Clone(): clone contains different header key: %+v",
					chain,
					clone.headerKey,
				)
			}

			if chain.headerKey != nil &&
				len(chain.headerKey.Bytes) > 0 &&
				&chain.headerKey.Bytes[0] == &clone.headerKey.Bytes[0] {
				t.Fatalf(
					"%+v.Clone(): clone contains same header key memory pointer %v",
					chain,
					clone.headerKey,
				)
			}

			if !reflect.DeepEqual(clone.nextHeaderKey, chain.nextHeaderKey) {
				t.Fatalf(
					"%+v.Clone(): clone contains different next header key: %+v",
					chain,
					clone.nextHeaderKey,
				)
			}

			if len(chain.nextHeaderKey.Bytes) > 0 &&
				&chain.nextHeaderKey.Bytes[0] == &clone.nextHeaderKey.Bytes[0] {
				t.Fatalf(
					"%+v.Clone(): clone contains same next header key key memory pointer %v",
					chain,
					clone.nextHeaderKey,
				)
			}
		})
	}
}

type chainEncryptTestArgs struct {
	masterKey       *keys.Master
	headerKey       *keys.Header
	nextHeaderKey   keys.Header
	headerPublicKey keys.Public
	data            []byte
	auth            []byte
}

var chainEncryptTests = []struct {
	name          string
	args          chainEncryptTestArgs
	errCategories []error
}{
	{
		"nil master key",
		chainEncryptTestArgs{
			nil,
			&keys.Header{
				Bytes: make([]byte, cipher.KeySize),
			},
			keys.Header{},
			keys.Public{},
			[]byte{6, 7, 8, 9},
			[]byte{3, 2, 1},
		},
		"advance chain: invalid value: master key is nil",
		[]error{
			ErrAdvance,
			ErrMasterKeyIsNil,
		},
	},
	{
		"nil header key",
		chainEncryptTestArgs{
			&keys.Master{
				Bytes: []byte{1, 2, 3},,
			},
			nil,
			keys.Header{},
			keys.Public{},
			[]byte{6, 7, 8, 9},
			[]byte{3, 2, 1},
		},
		[]error{
			ErrHeaderKeyIsNil,
		}
	},
	{
		"short header key",
		chainEncryptTestArgs{
			&keys.Master{
				Bytes: []byte{1, 2, 3},
			},
			&keys.Header{
				Bytes: []byte{1, 2, 3},
			},
			keys.Header{},
			keys.Public{},
			[]byte{6, 7, 8, 9},
			[]byte{3, 2, 1},
		},
		[]errors{
			ErrEncryptHeader,
			ErrNewCipher,
		},
	},
	{
		"success",
		chainEncryptTestArgs{
			&keys.Master{
				Bytes: []byte{1, 2, 3},
			},
			&keys.Header{
				Bytes: make([]byte, cipher.KeySize),
			},
			keys.Header{
				Bytes: make([]byte, cipher.KeySize),
			},
			keys.Public{
				Bytes: []byte{0, 1, 2, 3, 4, 5},
			},
			[]byte{6, 7, 8, 9},
			[]byte{3, 2, 1},
		},
		nil,
	},
}

func TestChainEncrypt(t *testing.T) {
	t.Parallel()

	for testIndex, test := range chainEncryptTests {
		testIndex := testIndex
		test := test

		chain, err := New(
			test.args.masterKey,
			test.args.headerKey,
			test.args.nextHeaderKey,
			uint64(testIndex),
			2,
		)
		if err != nil {
			t.Fatalf("New(): expected no error but got %v", err)
		}

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			header := chain.PrepareHeader(test.args.headerPublicKey)
			if header.MessageNumber != uint64(testIndex) {
				t.Fatalf(
					"expected header message number %d but got %d",
					testIndex,
					header.MessageNumber,
				)
			}

			encryptedHeader, encryptedData, err := chain.Encrypt(header, test.data, test.auth)
			if err != nil && len(test.errCategories) == 0 {
				t.Fatalf(
					"%+v.Encrypt(%+v, %v, %v) expected no error but got %v",
					chain,
					header,
					test.data,
					test.auth,
					err,
				)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, test.errCategory) {
					t.Fatalf(
						"%+v.Encrypt(%+v, %v, %v) expected error category %v but got %v",
						chain,
						header,
						test.data,
						test.auth,
						test.errCategory,
						err,
					)
				}
			}

			if err != nil {
				return
			}

			if len(encryptedHeader) == 0 {
				t.Fatalf(
					"%+v.Encrypt(%+v, %v, %v) returned empty encrypted header",
					chain,
					header,
					test.data,
					test.auth,
				)
			}

			if len(encryptedData) == 0 {
				t.Fatalf(
					"%+v.Encrypt(%+v, %v, %v) returned empty encrypted data",
					chain,
					header,
					test.data,
					test.auth,
				)
			}

			if reflect.DeepEqual(encryptedHeader, header.Encode()) {
				t.Fatalf(
					"%+v.Encrypt(%+v, %v, %v) returned input header bytes",
					chain,
					header,
					test.data,
					test.auth,
				)
			}

			if reflect.DeepEqual(encryptedData, test.data) {
				t.Fatalf(
					"%+v.Encrypt(%+v, %v, %v) returned input data",
					chain,
					header,
					test.data,
					test.auth,
				)
			}
		})
	}
}

type chainPrepareHeaderTestArgs struct {
	publicKey                  keys.Public
	nextMessageNumber          uint64
	previousChainMessagesCount uint64
}

var chainPrepareHeaderTests = []struct {
	name   string
	args   chainPrepareHeaderTestArgs
	header header.Header
}{
	{"zero args and header", chainPrepareHeaderTestArgs{}, header.Header{}},
	{
		"full args and header",
		chainPrepareHeaderTestArgs{
			keys.Public{
				Bytes: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			123,
			456,
		},
		header.Header{
			PublicKey: keys.Public{
				Bytes: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			MessageNumber:                     123,
			PreviousSendingChainMessagesCount: 456,
		},
	},
}

func TestChainPrepareHeader(t *testing.T) {
	t.Parallel()

	for _, test := range chainPrepareHeaderTests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chain, err := New(
				nil,
				nil,
				keys.Header{},
				test.args.nextMessageNumber,
				test.args.previousChainMessagesCount,
			)
			if err != nil {
				t.Fatalf("New(): expected no error but got %v", err)
			}

			header := chain.PrepareHeader(test.args.publicKey)
			if !reflect.DeepEqual(header, test.header) {
				t.Fatalf(
					"%+v.PrepareHeader(%+v): expected %+v but got %+v",
					chain,
					test.args.publicKey,
					test.header,
					header,
				)
			}
		})
	}
}

func TestChainUpgrade(t *testing.T) {
	t.Parallel()

	oldNextHeaderKey := keys.Header{
		Bytes: []byte{1, 2, 3},
	}

	oldNextMessageNumber := uint64(222)

	chain, err := New(nil, nil, oldNextHeaderKey, oldNextMessageNumber, 111)
	if err != nil {
		t.Fatalf("New(): expected no error but got %v", err)
	}

	masterKey := keys.Master{
		Bytes: []byte{11, 22, 33},
	}
	nextHeaderKey := keys.Header{
		Bytes: []byte{44, 55, 66, 77},
	}
	chain.Upgrade(masterKey, nextHeaderKey)

	if !reflect.DeepEqual(*chain.masterKey, masterKey) {
		t.Fatalf(
			"Upgrade(%+v, %+v): set different master key %+v",
			masterKey,
			nextHeaderKey,
			*chain.masterKey,
		)
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
		t.Fatalf(
			"Upgrade(%+v, %+v): set different next header key %+v",
			masterKey,
			nextHeaderKey,
			chain.nextHeaderKey,
		)
	}

	if chain.nextMessageNumber != 0 {
		t.Fatalf(
			"Upgrade(%+v, %+v): expected message number 0 but got %d",
			masterKey,
			nextHeaderKey,
			chain.nextMessageNumber,
		)
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
