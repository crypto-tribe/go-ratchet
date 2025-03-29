package header

import (
	"encoding/binary"
	"fmt"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
	"github.com/platform-inf/go-utils"
)

type Header struct {
	PublicKey                         keys.Public
	PreviousSendingChainMessagesCount uint64
	MessageNumber                     uint64
}

func Decode(bytes []byte) (Header, error) {
	if len(bytes) < 2*utils.Uint64Size {
		return Header{}, fmt.Errorf("%w: not enough bytes", errlist.ErrInvalidValue)
	}

	header := Header{}
	header.MessageNumber = binary.LittleEndian.Uint64(bytes[:utils.Uint64Size])
	header.PreviousSendingChainMessagesCount = binary.LittleEndian.Uint64(bytes[utils.Uint64Size : 2*utils.Uint64Size])

	if len(bytes) > 2*utils.Uint64Size {
		header.PublicKey = keys.Public{Bytes: bytes[2*utils.Uint64Size:]}
	}

	return header, nil
}

func (h Header) Encode() []byte {
	var messageNumberBytes, previousMessagesCountBytes [utils.Uint64Size]byte

	binary.LittleEndian.PutUint64(messageNumberBytes[:], h.MessageNumber)
	binary.LittleEndian.PutUint64(previousMessagesCountBytes[:], h.PreviousSendingChainMessagesCount)

	return utils.ConcatByteSlices(messageNumberBytes[:], previousMessagesCountBytes[:], h.PublicKey.Bytes)
}
