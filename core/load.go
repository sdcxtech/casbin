package core

// The iterator interface for loading policies and role mappings.
type LoadIterator interface {
	Next() (ok bool, key string, vals []string)
	Error() (err error)
}
