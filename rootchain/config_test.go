package rootchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-utils/check"
)

type testCrypto struct{}

func (tc testCrypto) AdvanceChain(
	_ keys.Root,
	_ keys.Shared,
) (keys.Root, keys.Master, keys.Header, error) {
	return keys.Root{}, keys.Master{}, keys.Header{}, nil
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
}

func TestNewConfigWithCryptoOption(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig(WithCrypto(testCrypto{}))
	if err != nil {
		t.Fatalf("newConfig() with options expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
		t.Fatal("WithCrypto() option did not set passed crypto")
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
