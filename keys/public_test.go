package keys

import (
	"reflect"
	"testing"
)

func TestPublicClone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  Public
	}{
		{"zero public key", Public{}},
		{"full public key", Public{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clone := test.key.Clone()
			if !reflect.DeepEqual(clone, test.key) {
				t.Fatalf("%+v.Clone() returned different value %+v", test.key, clone)
			}

			if len(clone.Bytes) > 0 && &clone.Bytes[0] == &test.key.Bytes[0] {
				t.Fatalf("%+v.Clone() returned same bytes memory %p", test.key, &clone.Bytes[0])
			}
		})
	}
}

func TestPublicClonePtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  *Public
	}{
		{"nil ptr to public key", nil},
		{"ptr to zero public key", &Public{}},
		{"ptr to full public key", &Public{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clone := test.key.ClonePtr()
			if !reflect.DeepEqual(clone, test.key) {
				t.Fatalf("%+v.ClonePtr() returned different value %+v", test.key, clone)
			}

			if clone != nil && len(clone.Bytes) > 0 && &clone.Bytes[0] == &test.key.Bytes[0] {
				t.Fatalf("%+v.Clone() returned same bytes memory %p", test.key, &clone.Bytes[0])
			}
		})
	}
}
