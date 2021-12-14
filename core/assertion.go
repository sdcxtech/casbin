package core

import (
	"fmt"
	"strings"
)

type AssertionSchema struct {
	indexToKey map[int]string
	keyToIndex map[string]int
}

func (s *AssertionSchema) Has(key string) bool {
	_, ok := s.keyToIndex[key]
	return ok
}

func (s AssertionSchema) CreateAssertion(vals []string) (assertion Assertion, err error) {
	expectCount := len(s.indexToKey)
	gotCount := len(vals)
	if expectCount != gotCount {
		err = fmt.Errorf("invalid values count: expect %d got %d", expectCount, gotCount)
		return
	}
	a := make(map[string]string, gotCount)
	for i, v := range vals {
		a[s.indexToKey[i]] = v
	}
	assertion = a
	return
}

func NewAssertionSchema(line string) (a AssertionSchema, err error) {
	subs := strings.Split(line, ",")
	if len(subs) < 3 {
		err = fmt.Errorf("%w: at least 3 column", ErrInvalidAssertionSchema)
		return
	}

	tokens := make([]string, 0, len(subs))
	for _, t := range subs {
		t = strings.TrimSpace(t)
		if len(t) == 0 {
			err = fmt.Errorf("%w: can't be empty", ErrInvalidAssertionSchema)
			return
		}
		tokens = append(tokens, t)
	}

	a = AssertionSchema{
		indexToKey: make(map[int]string, len(tokens)),
		keyToIndex: make(map[string]int, len(tokens)),
	}
	for i, t := range tokens {
		a.indexToKey[i] = t
		a.keyToIndex[t] = i
	}

	return
}

type AssertionIterator interface {
	Next() (ok bool, key string, vals []string)
	Error() (err error)
}

type Policies []Assertion

type Assertion map[string]string