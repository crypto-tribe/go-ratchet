package keys

import "testing"

func TestMessageMasterClone(t *testing.T) {
	t.Parallel()

	keys := []MessageMaster{
		MessageMaster{},
		MessageMaster{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.Clone()
		testBytesClone(t, key.Bytes, clone.Bytes)
	}
}

func TestMessageMasterClonePtr(t *testing.T) {
	t.Parallel()

	keys := []*MessageMaster{
		nil,
		&MessageMaster{},
		&MessageMaster{Bytes: []byte{1, 2, 3, 4, 5}},
	}

	for _, key := range keys {
		clone := key.ClonePtr()
		testClonePtrPointers(t, key, clone)

		if key != nil && clone != nil {
			testBytesClone(t, key.Bytes, clone.Bytes)
		}
	}
}
