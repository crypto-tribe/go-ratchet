package ratchet

import (
	"github.com/platform-source/aegis/keys"
)

// Crypto is the interface for rachet crypto.
type Crypto interface {
	ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error)
	GenerateKeyPair() (keys.Private, keys.Public, error)
}
