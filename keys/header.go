package keys

import "github.com/platform-inf/go-utils"

type Header struct {
	Bytes []byte
}

func (hk Header) Clone() Header {
	hk.Bytes = utils.CloneByteSlice(hk.Bytes)
	return hk
}

func (hk *Header) ClonePtr() *Header {
	if hk == nil {
		return nil
	}

	clone := hk.Clone()

	return &clone
}
