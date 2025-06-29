package maybe

import "errors"

var (
	// ErrMissingValue is returned when attempting to access a value that doesn't exist,
	// such as unwrapping a None Option or accessing an invalid Nullable value.
	ErrMissingValue = errors.New("expected a value, but none was present")
)
