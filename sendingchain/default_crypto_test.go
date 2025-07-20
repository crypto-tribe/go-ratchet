package sendingchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
	cipher "golang.org/x/crypto/chacha20poly1305"
)

var defaultCryptoAdvanceChainTests = []struct {
	name      string
	masterKey keys.Master
}{
	{
		"zero message master key",
		keys.Master{},
	},
	{
		"non-empty message master key",
		keys.Master{
			Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8},
		},
	},
}

func TestDefaultCryptoAdvanceChain(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	for _, test := range defaultCryptoAdvanceChainTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			masterKey, messageKey, err := crypto.AdvanceChain(test.masterKey)
			if err != nil {
				t.Fatalf("AdvanceChain(%+v): expected no error but got %v", test.masterKey, err)
			}

			if reflect.DeepEqual(masterKey, test.masterKey) {
				t.Fatalf("AdvanceChain(%+v): expected different message master key", test.masterKey)
			}

			if len(masterKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v): returned empty message master key", test.masterKey)
			}

			if len(messageKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v): returned empty message key", test.masterKey)
			}
		})
	}
}

var defaultCryptoEncryptHeaderTests = []struct {
	name          string
	headerKey     keys.Header
	header        header.Header
	errCategories []error
}{
	{
		"zero header key and zero header",
		keys.Header{},
		header.Header{},
		[]error{
			ErrNewCipher,
		},
	},
	{
		"invalid header key and zero header",
		keys.Header{
			Bytes: make([]byte, cipher.KeySize+1),
		},
		header.Header{},
		[]error{
			ErrNewCipher,
		},
	},
	{
		"non-empty header key and zero header",
		keys.Header{
			Bytes: make([]byte, cipher.KeySize),
		},
		header.Header{},
		nil,
	},
	{
		"non-empty header key and full header",
		keys.Header{
			Bytes: make([]byte, cipher.KeySize),
		},
		header.Header{
			PublicKey: keys.Public{
				Bytes: []byte{1, 2, 3, 4},
			},
			MessageNumber:                     222,
			PreviousSendingChainMessagesCount: 55,
		},
		nil,
	},
}

func TestDefaultCryptoEncryptHeader(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	for _, test := range defaultCryptoEncryptHeaderTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			encryptedHeader, err := crypto.EncryptHeader(test.headerKey, test.header)
			if err != nil && len(test.errCategories) == 0 {
				t.Fatalf(
					"EncryptHeader(%+v, %+v): expected no error but got %v",
					test.headerKey,
					test.header,
					err,
				)
			}

			for _, errCategory := range test.errCategories {
				if !errors.Is(err, errCategory) {
					t.Fatalf(
						"EncryptHeader(%+v, %+v): expected error %v but got %v",
						test.headerKey,
						test.header,
						errCategory,
						err,
					)
				}
			}

			if err != nil {
				return
			}

			if len(encryptedHeader) == 0 {
				t.Fatalf(
					"EncryptHeader(%+v, %+v): returned empty bytes",
					test.headerKey,
					test.header,
				)
			}
		})
	}
}

var defaultCryptoEncryptMessageTests = []struct {
	name       string
	messageKey keys.Message
	message    []byte
	auth       []byte
}{
	{
		"zero data",
		keys.Message{},
		nil,
		nil,
	},
	{
		"zero message key, without auth",
		keys.Message{},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		nil,
	},
	{
		"non-empty message key, without auth",
		keys.Message{
			Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		nil,
	},
	{
		"zero message key, with auth",
		keys.Message{},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]byte{1, 2, 3},
	},
	{
		"non-empty message key, with auth",
		keys.Message{
			Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]byte{1, 2, 3},
	},
}

func TestDefaultCryptoEncryptMessage(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	for _, test := range defaultCryptoEncryptMessageTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			encryptedMessage, err := crypto.EncryptMessage(test.messageKey, test.message, test.auth)
			if err != nil {
				t.Fatalf(
					"EncryptMessage(%+v, %+v, %+v): expected no error but got %v",
					test.messageKey,
					test.message,
					test.auth,
					err,
				)
			}

			if len(encryptedMessage) == 0 {
				t.Fatalf(
					"EncryptMessage(%+v, %+v, %+v): returned empty bytes",
					test.messageKey,
					test.message,
					test.auth,
				)
			}
		})
	}
}
