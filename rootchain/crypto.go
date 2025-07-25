package rootchain

import (
	"github.com/platform-source/aegis/keys"
)

// Crypto is the crypto interface for the root chain.
type Crypto interface {
	AdvanceChain(
		rootKey keys.Root,
		sharedKey keys.Shared,
	) (keys.Root, keys.Master, keys.Header, error)
}
