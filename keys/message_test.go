package keys

import "testing"

func TestMessageClone(t *testing.T) {
	t.Parallel()

	keys := []Message{
		Message{},
		Message{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.Clone()
		testBytesClone(t, key.Bytes, clone.Bytes)
	}
}
