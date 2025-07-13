package ratchet

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-ratchet/receivingchain"
	"github.com/lyreware/go-ratchet/rootchain"
	"github.com/lyreware/go-ratchet/sendingchain"
	"github.com/lyreware/go-utils/check"
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

func (tc testRootChainCrypto) AdvanceChain(
	_ keys.Root,
	_ keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	return keys.Root{}, keys.Master{}, keys.Header{}, nil
}

type testSendingChainCrypto struct{}

func (tc testSendingChainCrypto) AdvanceChain(_ keys.Master) (keys.Master, keys.Message, error) {
	return keys.Master{}, keys.Message{}, nil
}

func (tc testSendingChainCrypto) EncryptHeader(_ keys.Header, _ header.Header) ([]byte, error) {
	return nil, nil
}

func (tc testSendingChainCrypto) EncryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("default config", func(t *testing.T) {
		t.Parallel()

		cfg, err := newConfig()
		if err != nil {
			t.Fatalf("newConfig() expected no error but got %v", err)
		}

		if check.IsNil(cfg.crypto) {
			t.Fatal("newConfig() sets no default value for crypto")
		}
	})

	t.Run("chain options", func(t *testing.T) {
		t.Parallel()

		cfg, err := newConfig(
			WithReceivingChainOptions(receivingchain.WithCrypto(testReceivingChainCrypto{})),
			WithRootChainOptions(rootchain.WithCrypto(testRootChainCrypto{})),
			WithSendingChainOptions(sendingchain.WithCrypto(testSendingChainCrypto{})),
		)
		if err != nil {
			t.Fatalf("newConfig() with options expected no error but got %v", err)
		}

		if len(cfg.receivingOptions) != 1 {
			t.Fatal("newConfig() with receiving chain options did not set passed crypto")
		}

		if len(cfg.rootOptions) != 1 {
			t.Fatal("newConfig() with root chain options did not set passed crypto")
		}

		if len(cfg.sendingOptions) != 1 {
			t.Fatal("newConfig() with sending chain options did not set passed crypto")
		}
	})

	t.Run("crypto option success", func(t *testing.T) {
		t.Parallel()

		cfg, err := newConfig(WithCrypto(testCrypto{}))
		if err != nil {
			t.Fatalf("newConfig() with options expected no error but got %v", err)
		}

		if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
			t.Fatal("WithCrypto() option did not set passed crypto")
		}
	})

	t.Run("crypto option error", func(t *testing.T) {
		t.Parallel()

		_, err := newConfig(WithCrypto(nil))
		if err == nil || err.Error() != "option: invalid value: crypto is nil" {
			t.Fatalf("WithCrypto(nil) expected error but got %v", err)
		}

		if !errors.Is(err, errlist.ErrOption) {
			t.Fatalf("WithCrypto(nil) error is not option error but %v", err)
		}

		if !errors.Is(err, errlist.ErrInvalidValue) {
			t.Fatalf("WithCrypto(nil) error is not invalid value error but %v", err)
		}
	})
}
