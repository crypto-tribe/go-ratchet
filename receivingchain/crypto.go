package receivingchain

import (
	"github.com/crypto-tribe/go-ratchet/header"
	"github.com/crypto-tribe/go-ratchet/keys"
)

// Crypto is a crypto for the receiving chain.
type Crypto interface {
	AdvanceChain(masterKey keys.Master) (keys.Master, keys.Message, error)
	DecryptHeader(key keys.Header, encryptedHeader []byte) (header.Header, error)
	DecryptMessage(key keys.Message, encryptedMessage, auth []byte) ([]byte, error)
}
