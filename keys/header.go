package keys

import (
	"github.com/platform-source/tools/slices"
)

// Header is the key for header encryption and decryption.
type Header struct {
	Bytes []byte
}

// Clone clones header key.
func (hk Header) Clone() Header {
	hk.Bytes = slices.CloneBytes(hk.Bytes)

	return hk
}

// ClonePtr clones header key pointer.
func (hk *Header) ClonePtr() *Header {
	if hk == nil {
		return nil
	}

	clone := hk.Clone()

	return &clone
}
