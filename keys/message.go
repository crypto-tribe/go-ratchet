package keys

import (
	"github.com/lyreware/go-utils/slices"
)

type Message struct {
	Bytes []byte
}

func (mk Message) Clone() Message {
	mk.Bytes = slices.CloneBytes(mk.Bytes)
	return mk
}
