package header

import (
	"encoding/binary"

	"github.com/crypto-tribe/go-ratchet/keys"
	"github.com/crypto-tribe/go-utils/sizes"
	"github.com/crypto-tribe/go-utils/slices"
)

// Header is the message header.
type Header struct {
	PublicKey                         keys.Public
	PreviousSendingChainMessagesCount uint64
	MessageNumber                     uint64
}

// Decode decodes header bytes to the struct.
func Decode(headerBytes []byte) (Header, error) {
	if len(headerBytes) < 2*sizes.Uint64 {
		return Header{}, ErrNotEnoughBytes
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

// Encode encodes header struct to the bytes slice.
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
