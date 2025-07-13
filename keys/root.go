package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Root struct {
	Bytes []byte
}

func (rk Root) Clone() Root {
	rk.Bytes = slices.CloneBytes(rk.Bytes)
	return rk
}
