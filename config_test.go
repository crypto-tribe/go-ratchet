package ratchet

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
	"github.com/crypto-tribe/go-ratchet/receivingchain"
	"github.com/crypto-tribe/go-ratchet/rootchain"
	"github.com/crypto-tribe/go-ratchet/sendingchain"
)

type testCrypto struct{}

func (testCrypto) ComputeSharedKey(_ keys.Private, _ keys.Public) (keys.Shared, error) {
	return keys.Shared{}, nil
}

func (testCrypto) GenerateKeyPair() (keys.Private, keys.Public, error) {
	return keys.Private{}, keys.Public{}, nil
}

type testReceivingChainCrypto struct{}

func (testReceivingChainCrypto) AdvanceChain(_ keys.Master) (keys.Master, keys.Message, error) {
	return keys.Master{}, keys.Message{}, nil
}

func (testReceivingChainCrypto) DecryptHeader(_ keys.Header, _ []byte) (header.Header, error) {
	return header.Header{}, nil
}

func (testReceivingChainCrypto) DecryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

type testRootChainCrypto struct{}

func (testRootChainCrypto) AdvanceChain(
	_ keys.Root,
	_ keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	return keys.Root{}, keys.Master{}, keys.Header{}, nil
}

type testSendingChainCrypto struct{}

func (testSendingChainCrypto) AdvanceChain(_ keys.Master) (keys.Master, keys.Message, error) {
	return keys.Master{}, keys.Message{}, nil
}

func (testSendingChainCrypto) EncryptHeader(_ keys.Header, _ header.Header) ([]byte, error) {
	return nil, nil
}

func (testSendingChainCrypto) EncryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

var newConfigTests = []struct {
	name                        string
	options                     []Option
	errCategories               []error
	expectedCrypto              Crypto
	expectedReceivingOptionsLen int
	expectedRootOptionsLen      int
	expectedSendingOptionsLen   int
}{
	{
		"default",
		nil,
		nil,
		defaultCrypto{},
		0,
		0,
		0,
	},
	{
		"all options success",
		[]Option{
			WithCrypto(testCrypto{}),
			WithReceivingChainOptions(receivingchain.WithCrypto(testReceivingChainCrypto{})),
			WithRootChainOptions(rootchain.WithCrypto(testRootChainCrypto{})),
			WithSendingChainOptions(sendingchain.WithCrypto(testSendingChainCrypto{})),
		},
		nil,
		testCrypto{},
		1,
		1,
		1,
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
		0,
		0,
		0,
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

			if len(cfg.receivingOptions) != test.expectedReceivingOptionsLen {
				t.Fatal("WithReceivingChainOptions() option did not set passed options")
			}

			if len(cfg.rootOptions) != test.expectedRootOptionsLen {
				t.Fatal("WithRootChainOptions() option did not set passed options")
			}

			if len(cfg.sendingOptions) != test.expectedSendingOptionsLen {
				t.Fatal("WithSendingChainOptions() option did not set passed options")
			}
		})
	}
}
