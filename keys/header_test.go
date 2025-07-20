package keys

import (
	"reflect"
	"testing"
)

var headerCloneTests = []struct {
	name string
	key  Header
}{
	{
		"zero header key",
		Header{},
	},
	{
		"non-empty header key",
		Header{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestHeaderClone(t *testing.T) {
	t.Parallel()

	for _, test := range headerCloneTests {
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

var headerClonePtrTests = []struct {
	name string
	key  *Header
}{
	{
		"nil ptr to header key",
		nil,
	},
	{
		"ptr to zero header key",
		&Header{},
	},
	{
		"ptr to non-empty header key",
		&Header{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestHeaderClonePtr(t *testing.T) {
	t.Parallel()

	for _, test := range headerClonePtrTests {
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
