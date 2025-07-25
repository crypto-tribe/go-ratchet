package keys

import (
	"github.com/platform-source/tools/slices"
)

// Public is participant's public key.
type Public struct {
	Bytes []byte
}

// Clone clones public key.
func (pk Public) Clone() Public {
	pk.Bytes = slices.CloneBytes(pk.Bytes)

	return pk
}

// ClonePtr clones public key pointer.
func (pk *Public) ClonePtr() *Public {
	if pk == nil {
		return nil
	}

	clone := pk.Clone()

	return &clone
}
