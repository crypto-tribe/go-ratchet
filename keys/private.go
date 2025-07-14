package keys

import (
	"github.com/lyreware/go-utils/slices"
)

// Private key is participant's private key.
type Private struct {
	Bytes []byte
}

// Clone clones private key.
func (pk Private) Clone() Private {
	pk.Bytes = slices.CloneBytes(pk.Bytes)

	return pk
}
