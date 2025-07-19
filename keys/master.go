package keys

import (
	"github.com/crypto-tribe/go-utils/slices"
)

// Master is the master key to derive new message keys.
type Master struct {
	Bytes []byte
}

// Clone clones master key.
func (mk Master) Clone() Master {
	mk.Bytes = slices.CloneBytes(mk.Bytes)

	return mk
}

// ClonePtr clones master key pointer.
func (mk *Master) ClonePtr() *Master {
	if mk == nil {
		return nil
	}

	clone := mk.Clone()

	return &clone
}
