package header

import (
	"encoding/binary"
	"fmt"

	"github.com/lyreware/go-ratchet/errlist"
	"github.com/lyreware/go-ratchet/keys"
	"github.com/lyreware/go-utils/sizes"
	"github.com/lyreware/go-utils/slices"
)

type Header struct {
	PublicKey                         keys.Public
	PreviousSendingChainMessagesCount uint64
	MessageNumber                     uint64
}

func Decode(headerBytes []byte) (Header, error) {
	if len(headerBytes) < 2*sizes.Uint64 {
		return Header{}, fmt.Errorf("%w: not enough bytes", errlist.ErrInvalidValue)
	}

	header := Header{
		MessageNumber: binary.LittleEndian.Uint64(
			headerBytes[:sizes.Uint64],
		),
		PreviousSendingChainMessagesCount: binary.LittleEndian.Uint64(
			headerBytes[sizes.Uint64 : 2*sizes.Uint64],
		),
	}

	if len(headerBytes) > 2*sizes.Uint64 {
		header.PublicKey = keys.Public{
			Bytes: headerBytes[2*sizes.Uint64:],
		}
	}

	return header, nil
}

func (h Header) Encode() []byte {
	var messageNumberBytes [sizes.Uint64]byte
	binary.LittleEndian.PutUint64(
		messageNumberBytes[:],
		h.MessageNumber,
	)

	var previousMessagesCountBytes [sizes.Uint64]byte
	binary.LittleEndian.PutUint64(
		previousMessagesCountBytes[:],
		h.PreviousSendingChainMessagesCount,
	)

	headerBytes := slices.ConcatBytes(
		messageNumberBytes[:],
		previousMessagesCountBytes[:],
		h.PublicKey.Bytes,
	)

	return headerBytes
}
