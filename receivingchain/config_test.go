package receivingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/header"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-utils/check"
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

func TestNewDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig()
	if err != nil {
		t.Fatalf("newConfig() expected no error but got %v", err)
	}

	if check.IsNil(cfg.crypto) {
		t.Fatal("newConfig() sets no default value for crypto")
	}

	if check.IsNil(cfg.skippedKeysStorage) {
		t.Fatal("newConfig() sets no default value for skipped keys storage")
	}
}

func TestNewConfigWithCryptoOptionSuccess(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig(WithCrypto(testCrypto{}))
	if err != nil {
		t.Fatalf("newConfig() with crypto option expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
		t.Fatal("WithCrypto() option did not set passed crypto")
	}
}

func TestNewConfigWithSkippedKeysOptionSuccess(t *testing.T) {
	t.Parallel()

	var storage testSkippedKeysStorage

	cfg, err := newConfig(WithSkippedKeysStorage(&storage))
	if err != nil {
		t.Fatalf("newConfig() with skipped keys storage options expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.skippedKeysStorage) != reflect.TypeOf(&storage) {
		t.Fatal("WithSkippedKeysStorage() option did not set passed skipped keys storage")
	}
}

func TestNewConfigWithCryptoOptionError(t *testing.T) {
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
}

func TestNewConfigWithSkippedKeysOptionError(t *testing.T) {
	t.Parallel()

	_, err := newConfig(WithSkippedKeysStorage(nil))
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
