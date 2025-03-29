package keys

import "github.com/platform-inf/go-utils"

type MessageMaster struct {
	Bytes []byte
}

func (mk MessageMaster) Clone() MessageMaster {
	mk.Bytes = utils.CloneByteSlice(mk.Bytes)
	return mk
}

func (mk *MessageMaster) ClonePtr() *MessageMaster {
	if mk == nil {
		return nil
	}

	clone := mk.Clone()

	return &clone
}
