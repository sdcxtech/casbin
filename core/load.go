package core

// The iterator interface for loading policies and role mappings.
type LoadIterator interface {
	// Next returns next policy or role mapping.
	//
	// The `key` may be `p` or a role mapping key like `g`.
	// If `ok` is false. There is no more data or an error.
	Next() (ok bool, key string, vals []string)

	// Error returns the error that has occurred in `Next()` calling.
	Error() (err error)
}
