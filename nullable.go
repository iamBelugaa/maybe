package maybe

// Nullable[T] represents a value that might be null.
// Unlike Option, Nullable is specifically designed for handling
// null values in external systems like databases and JSON APIs.
// This is particularly useful for database operations where NULL
// values are common and need distinct handling.
type Nullable[T any] struct {
	value T
	valid bool
}
