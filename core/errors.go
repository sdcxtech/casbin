package core

import "errors"

var (
	ErrNotFoundMatcher        = errors.New("not found specified matcher")
	ErrAtLeastOneMatcher      = errors.New("should have at least one matcher")
	ErrInvalidAssertionSchema = errors.New("invalid assertion schema")
)
