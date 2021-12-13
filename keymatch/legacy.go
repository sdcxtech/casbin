package keymatch

import (
	"strings"
)

// KeyMatch determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*"
func KeyMatch(key1 string, key2 string) (bool, error) {
	i := strings.Index(key2, "*")
	if i == -1 {
		return key1 == key2, nil
	}

	if len(key1) > i {
		return key1[:i] == key2[:i], nil
	}
	return key1 == key2[:i], nil
}
