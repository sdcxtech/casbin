package keymatch

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/sdcxtech/casbin/core"
)

func NewRegexMatch() core.KeyMatchFunc {
	memorized := &sync.Map{}

	return func(key1 string, pattern string) (matched bool, err error) {
		var re *regexp.Regexp
		v, ok := memorized.Load(pattern)
		if ok {
			re = v.(*regexp.Regexp)
		} else {
			re, err = regexp.Compile(pattern)
			if err != nil {
				err = fmt.Errorf("compile regex: %s", pattern)
				return
			}
			memorized.Store(pattern, re)
		}

		matched = re.MatchString(key1)
		return
	}
}
