package keys

import (
	"reflect"
	"testing"
)

var messageMasterCloneTests = []struct {
	name string
	key  Master
}{
	{"zero message master key", Master{}},
	{
		"non-empty message master key",
		Master{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestMasterClone(t *testing.T) {
	t.Parallel()

	for _, test := range messageMasterCloneTests {
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

var messageMasterClonePtrTests = []struct {
	name string
	key  *Master
}{
	{"nil ptr to message master key", nil},
	{"ptr to zero message master key", &Master{}},
	{"ptr to non-empty message master key", &Master{Bytes: []byte{1, 2, 3, 4, 5}}},
}

func TestMasterClonePtr(t *testing.T) {
	t.Parallel()

	for _, test := range messageMasterClonePtrTests {
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
