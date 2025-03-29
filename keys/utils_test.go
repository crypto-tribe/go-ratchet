package keys

import (
	"slices"
	"testing"
)

func testBytesClone(t *testing.T, bytes, clone []byte) {
	if !slices.Equal(clone, bytes) {
		t.Fatalf("different key bytes: %v != %v", bytes, clone)
	}

	if len(bytes) > 0 && &bytes[0] == &clone[0] {
		t.Fatalf("the same slices: %v is %v", bytes, clone)
	}
}

func testClonePtrPointers[V any](t *testing.T, original, clone *V) {
	// Both are nil or not nil different memory pointers.
	if (original != nil || clone != nil) && (original == nil || clone == nil || original == clone) {
		t.Fatalf("%+v.ClonePtr() returned invalid memory %+v", original, clone)
	}
}
