package keymatch

import (
	"fmt"
	"path"
)

// GlobMatch determines whether key1 matches the pattern of key2 using glob pattern.
func GlobMatch(key1 string, key2 string) (bool, error) {
	matched, err := path.Match(key2, key1)
	if err != nil {
		err = fmt.Errorf("path match: %w", err)

		return matched, err
	}

	return matched, err
}
