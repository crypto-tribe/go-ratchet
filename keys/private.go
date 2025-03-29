package keys

import "github.com/platform-inf/go-utils"

type Private struct {
	Bytes []byte
}

func (pk Private) Clone() Private {
	pk.Bytes = utils.CloneByteSlice(pk.Bytes)
	return pk
}
