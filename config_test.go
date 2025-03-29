package ratchet

import (
	"errors"
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-ratchet/receivingchain"
	"github.com/platform-inf/go-ratchet/rootchain"
	"github.com/platform-inf/go-ratchet/sendingchain"
	"github.com/platform-inf/go-utils"
)

type testCrypto struct{}

func (tc testCrypto) ComputeSharedKey(_ keys.Private, _ keys.Public) (keys.Shared, error) {
	return keys.Shared{}, nil
}

func (tc testCrypto) GenerateKeyPair() (keys.Private, keys.Public, error) {
	return keys.Private{}, keys.Public{}, nil
}

type testReceivingChainCrypto struct{}

func (tc testReceivingChainCrypto) AdvanceChain(_ keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
	return keys.MessageMaster{}, keys.Message{}, nil
}

func (tc testReceivingChainCrypto) DecryptHeader(_ keys.Header, _ []byte) (header.Header, error) {
	return header.Header{}, nil
}

func (tc testReceivingChainCrypto) DecryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

type testRootChainCrypto struct{}

func (tc testRootChainCrypto) AdvanceChain(
	_ keys.Root,
	_ keys.Shared,
) (keys.Root, keys.MessageMaster, keys.Header, error) {
	return keys.Root{}, keys.MessageMaster{}, keys.Header{}, nil
}

type testSendingChainCrypto struct{}

func (tc testSendingChainCrypto) AdvanceChain(_ keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
	return keys.MessageMaster{}, keys.Message{}, nil
}

func (tc testSendingChainCrypto) EncryptHeader(_ keys.Header, _ header.Header) ([]byte, error) {
	return nil, nil
}

func (tc testSendingChainCrypto) EncryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

func TestNewConfigDefault(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig(nil)
	if err != nil {
		t.Fatalf("newConfig() expected no error but got %v", err)
	}

	if utils.IsNil(cfg.crypto) {
		t.Fatal("newConfig() sets no default value for crypto")
	}
}

func TestNewConfigWithChainOptions(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig([]Option{
		WithReceivingChainOptions(receivingchain.WithCrypto(testReceivingChainCrypto{})),
		WithRootChainOptions(rootchain.WithCrypto(testRootChainCrypto{})),
		WithSendingChainOptions(sendingchain.WithCrypto(testSendingChainCrypto{})),
	})
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
}

func TestNewConfigWithCryptoSuccess(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig([]Option{WithCrypto(testCrypto{})})
	if err != nil {
		t.Fatalf("newConfig() with options expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
		t.Fatal("WithCrypto() option did not set passed crypto")
	}
}

func TestNewConfigWithCryptoError(t *testing.T) {
	t.Parallel()

	_, err := newConfig([]Option{WithCrypto(nil)})
	if err == nil || err.Error() != "option: invalid value: crypto is nil" {
		t.Fatalf("WithCrypto(nil) expected error but got %v", err)
	}

	if !errors.Is(err, errlist.ErrOption) {
		t.Fatalf("WithCrypto(nil) error is not option error but %v", err)
	}

	if !errors.Is(err, errlist.ErrInvalidValue) {
		t.Fatalf("WithCrypto(nil) error is not invalid value error but %v", err)
	}
}
