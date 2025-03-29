package keys

import "testing"

func TestHeaderClone(t *testing.T) {
	t.Parallel()

	keys := []Header{
		Header{},
		Header{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.Clone()
		testBytesClone(t, key.Bytes, clone.Bytes)
	}
}

func TestHeaderClonePtr(t *testing.T) {
	t.Parallel()

	keys := []*Header{
		nil,
		&Header{},
		&Header{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.ClonePtr()
		testClonePtrPointers(t, key, clone)

		if key != nil && clone != nil {
			testBytesClone(t, key.Bytes, clone.Bytes)
		}
	}
}
