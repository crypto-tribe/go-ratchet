package receivingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-utils"
)

type testCrypto struct{}

func (tc testCrypto) AdvanceChain(_ keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
	return keys.MessageMaster{}, keys.Message{}, nil
}

func (tc testCrypto) DecryptHeader(_ keys.Header, _ []byte) (header.Header, error) {
	return header.Header{}, nil
}

func (tc testCrypto) DecryptMessage(_ keys.Message, _, _ []byte) ([]byte, error) {
	return nil, nil
}

type testSkippedKeysStorage struct {
	cloneCalled bool
}

func (ts *testSkippedKeysStorage) Add(_ keys.Header, _ uint64, _ keys.Message) error {
	return nil
}

func (ts *testSkippedKeysStorage) Clone() SkippedKeysStorage {
	ts.cloneCalled = true
	return ts
}

func (ts *testSkippedKeysStorage) Delete(_ keys.Header, _ uint64) error {
	return nil
}

func (ts *testSkippedKeysStorage) GetIter() (SkippedKeysIter, error) {
	return func(_ SkippedKeysYield) {}, nil
}

func TestConfigClone(t *testing.T) {
	t.Parallel()

	var skippedKeysStorage testSkippedKeysStorage

	cfg, err := newConfig([]Option{WithSkippedKeysStorage(&skippedKeysStorage)})
	if err != nil {
		t.Fatalf("newConfig() expected no error but got %v", err)
	}

	cfg.clone()

	if !skippedKeysStorage.cloneCalled {
		t.Fatal("clone() expected skipped keys storage clone")
	}
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

	if utils.IsNil(cfg.skippedKeysStorage) {
		t.Fatal("newConfig() sets no default value for skipped keys storage")
	}
}

func TestNewConfigWithCryptoSuccess(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig([]Option{WithCrypto(testCrypto{})})
	if err != nil {
		t.Fatalf("newConfig() with crypto option expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
		t.Fatal("WithCrypto() option did not set passed crypto")
	}
}

func TestNewConfigWithSkippedKeysStorageSuccess(t *testing.T) {
	t.Parallel()

	var storage testSkippedKeysStorage

	cfg, err := newConfig([]Option{WithSkippedKeysStorage(&storage)})
	if err != nil {
		t.Fatalf("newConfig() with skipped keys storage options expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.skippedKeysStorage) != reflect.TypeOf(&storage) {
		t.Fatal("WithSkippedKeysStorage() option did not set passed skipped keys storage")
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

func TestNewConfigWithSkippedKeysStorageError(t *testing.T) {
	t.Parallel()

	_, err := newConfig([]Option{WithSkippedKeysStorage(nil)})
	if err == nil || err.Error() != "option: invalid value: storage is nil" {
		t.Fatalf("WithSkippedKeysStorage(nil) expected error but got %v", err)
	}

	if !errors.Is(err, errlist.ErrOption) {
		t.Fatalf("WithSkippedKeysStorage(nil) error is not option error but %v", err)
	}

	if !errors.Is(err, errlist.ErrInvalidValue) {
		t.Fatalf("WithSkippedKeysStorage(nil) error is not invalid value error but %v", err)
	}
}
