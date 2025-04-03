package keys

import (
	"reflect"
	"testing"
)

func TestHeaderClone(t *testing.T) {
	t.Parallel()

	tests := []struct{
		name string
		key  Header
	}{
		{"zero header key", Header{}},
		{"full header key", Header{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clone := test.key.Clone()
			if !reflect.DeepEqual(clone, test.key) {
				t.Fatalf("%+v.Clone() returned different value %+v", test.key, clone)
			}
		})
	}
}

func TestHeaderClonePtr(t *testing.T) {
	t.Parallel()

	tests := []struct{
		name string
		key  *Header
	}{
		{"nil ptr to header key", nil},
		{"ptr to zero header key", &Header{}},
		{"ptr to full header key", &Header{Bytes: []byte{1, 2, 3, 4, 5}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clone := test.key.ClonePtr()
			if (test.key != nil || clone != nil) && (test.key == nil || clone == nil || test.key == clone) {
				t.Fatalf("%+v.ClonePtr() expected pointer %p but got %p", test.key, test.key, clone)
			}

			if !reflect.DeepEqual(clone, test.key) {
				t.Fatalf("%+v.ClonePtr() returned different value %+v", test.key, clone)
			}
		})
	}
}
