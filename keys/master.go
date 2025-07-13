package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Master struct {
	Bytes []byte
}

func (mk Master) Clone() Master {
	mk.Bytes = slices.CloneBytes(mk.Bytes)
	return mk
}

func (mk *Master) ClonePtr() *Master {
	if mk == nil {
		return nil
	}

	clone := mk.Clone()

	return &clone
}
