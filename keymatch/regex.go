package keymatch

import (
	"fmt"
	"regexp"
	"sync"
)

// NewRegexMatch build a regex match function.
func NewRegexMatch() Func {
	memorized := &sync.Map{}

	return func(key1 string, pattern string) (matched bool, err error) {
		var re *regexp.Regexp

		if v, ok := memorized.Load(pattern); ok {
			re = v.(*regexp.Regexp) // nolint: forcetypeassert
		} else {
			re, err = regexp.Compile(pattern)
			if err != nil {
				err = fmt.Errorf("compile regex: %w", err)

				return
			}
			memorized.Store(pattern, re)
		}

		matched = re.MatchString(key1)

		return
	}
}
