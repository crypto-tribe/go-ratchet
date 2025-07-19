package receivingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
)

type testCrypto struct{}

func (testCrypto) AdvanceChain(_ keys.Master) (keys.Master, keys.Message, error) {
	return keys.Master{}, keys.Message{}, nil
}

func (testCrypto) DecryptHeader(_ keys.Header, _ []byte) (header.Header, error) {
	return header.Header{}, nil
}

func (testCrypto) DecryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

type testSkippedKeysStorage struct {
	cloneCalled bool
}

func (*testSkippedKeysStorage) Add(_ keys.Header, _ uint64, _ keys.Message) error {
	return nil
}

func (ts *testSkippedKeysStorage) Clone() SkippedKeysStorage {
	ts.cloneCalled = true

	return ts
}

func (*testSkippedKeysStorage) Delete(_ keys.Header, _ uint64) error {
	return nil
}

func (*testSkippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	return func(_ SkippedKeysYield) {}, nil
}

var newConfigTests = []struct {
	name                       string
	options                    []Option
	errCategories              []error
	expectedCrypto             Crypto
	expectedSkippedKeysStorage SkippedKeysStorage
}{
	{
		"default",
		nil,
		nil,
		defaultCrypto{},
		defaultSkippedKeysStorage{},
	},
	{
		"crypto and skipped keys storage options success",
		[]Option{
			WithCrypto(testCrypto{}),
			WithSkippedKeysStorage(&testSkippedKeysStorage{}),
		},
		nil,
		testCrypto{},
		&testSkippedKeysStorage{},
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
		nil,
	},
	{
		"nil skipped keys storage",
		[]Option{
			WithSkippedKeysStorage(nil),
		},
		[]error{
			ErrApplyOptions,
			ErrSkippedKeysStorageIsNil,
		},
		nil,
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

			if reflect.TypeOf(
				cfg.skippedKeysStorage,
			) != reflect.TypeOf(
				test.expectedSkippedKeysStorage,
			) {
				t.Fatal("WithSkippedKeysStorage() option did not set passed skipped keys storage")
			}
		})
	}
}

func TestConfigClone(t *testing.T) {
	t.Parallel()

	var skippedKeysStorage testSkippedKeysStorage

	cfg, err := newConfig(WithSkippedKeysStorage(&skippedKeysStorage))
	if err != nil {
		t.Fatalf("newConfig() expected no error but got %v", err)
	}

	clone := cfg.clone()
	if !reflect.DeepEqual(clone, cfg) {
		t.Fatalf("%+v.clone() returned different config %+v", cfg, clone)
	}

	if !skippedKeysStorage.cloneCalled {
		t.Fatal("clone() expected skipped keys storage clone")
	}
}
