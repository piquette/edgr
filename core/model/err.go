package model

import "errors"

// Common errors.
var (
	ErrDuplicate         = errors.New("duplicate entry")
	ErrNotFound          = errors.New("not found")
	ErrDocumentsMismatch = errors.New("documents did not match count")
)
