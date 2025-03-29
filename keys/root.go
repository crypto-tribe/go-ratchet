package keys

import "github.com/platform-inf/go-utils"

type Root struct {
	Bytes []byte
}

func (rk Root) Clone() Root {
	rk.Bytes = utils.CloneByteSlice(rk.Bytes)
	return rk
}
