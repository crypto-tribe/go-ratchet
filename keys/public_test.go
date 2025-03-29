package keys

import "testing"

func TestPublicClone(t *testing.T) {
	t.Parallel()

	keys := []Public{
		Public{},
		Public{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.Clone()
		testBytesClone(t, key.Bytes, clone.Bytes)
	}
}

func TestPublicClonePtr(t *testing.T) {
	t.Parallel()

	keys := []*Public{
		nil,
		&Public{},
		&Public{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.ClonePtr()
		testClonePtrPointers(t, key, clone)

		if key != nil && clone != nil {
			testBytesClone(t, key.Bytes, clone.Bytes)
		}
	}
}
