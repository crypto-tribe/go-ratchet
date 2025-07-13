package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Public struct {
	Bytes []byte
}

func (pk Public) Clone() Public {
	pk.Bytes = slices.CloneBytes(pk.Bytes)

	return pk
}

func (pk *Public) ClonePtr() *Public {
	if pk == nil {
		return nil
	}

	clone := pk.Clone()

	return &clone
}
