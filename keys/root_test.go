package keys

import "testing"

func TestRootClone(t *testing.T) {
	t.Parallel()

	keys := []Root{
		Root{},
		Root{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.Clone()
		testBytesClone(t, key.Bytes, clone.Bytes)
	}
}
