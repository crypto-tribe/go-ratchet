package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Private struct {
	Bytes []byte
}

func (pk Private) Clone() Private {
	pk.Bytes = slices.CloneBytes(pk.Bytes)
	return pk
}
