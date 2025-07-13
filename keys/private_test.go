package keys

import (
	"reflect"
	"testing"
)

var privateCloneTests = []struct {
	name string
	key  Private
}{
	{"zero private key", Private{}},
	{
		"non-empty private key",
		Private{
			Bytes: []byte{1, 2, 3, 4, 5},
		},
	},
}

func TestPrivateClone(t *testing.T) {
	t.Parallel()

	for _, test := range privateCloneTests {
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
