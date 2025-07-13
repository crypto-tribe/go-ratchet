package keys

import (
	"reflect"
	"testing"
)

var rootCloneTests = []struct {
	name string
	key  Root
}{
	{"zero root key", Root{}},
	{
		"non-empty root key",
		Root{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestRootClone(t *testing.T) {
	t.Parallel()

	for _, test := range rootCloneTests {
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
