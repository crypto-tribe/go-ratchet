package keys

import "github.com/platform-inf/go-utils"

type Message struct {
	Bytes []byte
}

func (mk Message) Clone() Message {
	mk.Bytes = utils.CloneByteSlice(mk.Bytes)
	return mk
}
