package sendingchain

import (
	"reflect"
	"testing"

	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/platform-inf/go-ratchet/header"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestDefaultCryptoAdvanceChain(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	tests := []struct {
		name             string
		messageMasterKey keys.MessageMaster
	}{
		{"zero message master key", keys.MessageMaster{}},
		{"full message master key", keys.MessageMaster{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			messageMasterKey, messageKey, err := crypto.AdvanceChain(test.messageMasterKey)
			if err != nil {
				t.Fatalf("AdvanceChain(%+v): expected no error but got %v", test.messageMasterKey, err)
			}

			if reflect.DeepEqual(messageMasterKey, test.messageMasterKey) {
				t.Fatalf("AdvanceChain(%+v): expected different message master key", test.messageMasterKey)
			}

			if len(messageMasterKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v): returned empty message master key", test.messageMasterKey)
			}

			if len(messageKey.Bytes) == 0 {
				t.Fatalf("AdvanceChain(%+v): returned empty message key", test.messageMasterKey)
			}
		})
	}
}

func TestDefaultCryptoEncryptHeader(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	tests := []struct {
		name      string
		headerKey keys.Header
		header    header.Header
		errString string
	}{
		{
			"zero header key and zero header",
			keys.Header{},
			header.Header{},
			"new cipher: chacha20poly1305: bad key length",
		},
		{
			"invalid header key and zero header",
			keys.Header{Bytes: make([]byte, cipher.KeySize+1)},
			header.Header{},
			"new cipher: chacha20poly1305: bad key length",
		},
		{
			"full header key and zero header",
			keys.Header{Bytes: make([]byte, cipher.KeySize)},
			header.Header{},
			"",
		},
		{
			"full header key and full header",
			keys.Header{Bytes: make([]byte, cipher.KeySize)},
			header.Header{
				PublicKey:                         keys.Public{Bytes: []byte{1, 2, 3, 4}},
				MessageNumber:                     222,
				PreviousSendingChainMessagesCount: 55,
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			encryptedHeader, err := crypto.EncryptHeader(test.headerKey, test.header)
			if (err == nil && test.errString != "") || (err != nil && err.Error() != test.errString) {
				t.Fatalf("EncryptHeader(%+v, %+v): expected err %q but got %+v", test.headerKey, test.header, test.errString, err)
			}

			if err != nil {
				return
			}

			if len(encryptedHeader) == 0 {
				t.Fatalf("EncryptHeader(%+v, %+v): returned empty bytes", test.headerKey, test.header)
			}
		})
	}
}

func TestDefaultCryptoEncryptMessage(t *testing.T) {
	t.Parallel()

	crypto := newDefaultCrypto()

	tests := []struct {
		name       string
		messageKey keys.Message
		message    []byte
		auth       []byte
	}{
		{"zero data", keys.Message{}, nil, nil},
		{
			"zero message key, without auth",
			keys.Message{},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			nil,
		},
		{
			"full message key, without auth",
			keys.Message{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
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
			"full message key, with auth",
			keys.Message{Bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]byte{1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			encryptedMessage, err := crypto.EncryptMessage(test.messageKey, test.message, test.auth)
			if err != nil {
				t.Fatalf(
					"EncryptMessage(%+v, %+v, %+v): expected no error but got %v", test.messageKey, test.message, test.auth, err)
			}

			if len(encryptedMessage) == 0 {
				t.Fatalf("EncryptMessage(%+v, %+v, %+v): returned empty bytes", test.messageKey, test.message, test.auth)
			}
		})
	}
}
