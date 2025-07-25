package keys

import (
	"github.com/platform-source/tools/slices"
)

// Root is the key of ratchet root chain.
type Root struct {
	Bytes []byte
}

// Clone clones root key.
func (rk Root) Clone() Root {
	rk.Bytes = slices.CloneBytes(rk.Bytes)

	return rk
}
