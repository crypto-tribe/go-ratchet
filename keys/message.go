package keys

import (
	"github.com/platform-source/tools/slices"
)

// Message is the key to encrypt or decrypt messages.
type Message struct {
	Bytes []byte
}

// Clone clones message key.
func (mk Message) Clone() Message {
	mk.Bytes = slices.CloneBytes(mk.Bytes)

	return mk
}
