package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Header struct {
	Bytes []byte
}

func (hk Header) Clone() Header {
	hk.Bytes = slices.CloneBytes(hk.Bytes)
	return hk
}

func (hk *Header) ClonePtr() *Header {
	if hk == nil {
		return nil
	}

	clone := hk.Clone()

	return &clone
}
