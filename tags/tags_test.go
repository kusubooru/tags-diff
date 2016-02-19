package tags_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kusubooru/tags-diff/tags"
)

var diffTests = []struct {
	old     []string
	new     []string
	removed []string
	added   []string
}{

	{[]string{}, []string{"a"}, []string{}, []string{"a"}},
	{[]string{"a"}, []string{}, []string{"a"}, []string{}},
	{[]string{"a", "b"}, []string{"b", "c"}, []string{"a"}, []string{"c"}},
	{[]string{"a", "a", "b"}, []string{"b", "b", "c"}, []string{"a", "a"}, []string{"c"}},
	{[]string{"a:a", "b:b"}, []string{"b:b", "c:c"}, []string{"a:a"}, []string{"c:c"}},
}

func TestDiff(t *testing.T) {
	for _, tt := range diffTests {
		r, a := tags.Diff(tt.old, tt.new)
		if !reflect.DeepEqual(r, tt.removed) || !reflect.DeepEqual(a, tt.added) {
			t.Errorf("tags.Diff(%q, %q) = %#v, %#v, want %#v, %#v", tt.old, tt.new, r, a, tt.removed, tt.added)
		}
	}
}

func TestDiffFields(t *testing.T) {
	for _, tt := range diffTests {
		r, a := tags.DiffFields(strings.Join(tt.old, " "), strings.Join(tt.new, " "))
		if !reflect.DeepEqual(r, tt.removed) || !reflect.DeepEqual(a, tt.added) {
			t.Errorf("tags.Diff(%q, %q) = %#v, %#v, want %#v, %#v", tt.old, tt.new, r, a, tt.removed, tt.added)
		}
	}
}
