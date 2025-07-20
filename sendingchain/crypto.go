package sendingchain

import (
	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
)

// Crypto interface for sending chain.
type Crypto interface {
	AdvanceChain(masterKey keys.Master) (keys.Master, keys.Message, error)
	EncryptHeader(key keys.Header, head header.Header) ([]byte, error)
	EncryptMessage(key keys.Message, message, auth []byte) ([]byte, error)
}
