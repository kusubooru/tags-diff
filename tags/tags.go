package tags

import "strings"

// DiffFields splits the strings old and new around each instance of one or
// more consecutive white space characters turning them into two slices of
// strings then finds the difference. It is equivalent to calling
// strings.Fields on old and new then Diff.
func DiffFields(old, new string) (removed, added []string) {
	a := strings.Fields(old)
	b := strings.Fields(new)
	return diff(a, b), diff(b, a)
}

// Diff finds the difference between two slices of strings old and new,
// returning a slice of strings that were removed from old and a slice of
// strings that were added to new.
func Diff(old, new []string) (removed, added []string) {
	return diff(old, new), diff(new, old)
}

func diff(a, b []string) []string {
	diff := []string{}

	for _, aa := range a {
		found := false
		for _, bb := range b {
			if aa == bb {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, aa)
		}
	}
	return diff
}
