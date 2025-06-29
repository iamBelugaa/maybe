package maybe

// Option[T] represents an optional value: either Some(value) or None.
// It provides a type-safe alternative to nil pointers and helps avoid nil pointer panics.
// The zero value of Option is None (no value present).
type Option[T any] struct {
	value T
	has   bool
}

// Some creates a new Option containing the provided value.
// Use this when you have a definite value to wrap.
func Some[T any](value T) Option[T] {
	return Option[T]{has: true, value: value}
}

// None creates a new Option with no value.
// Use this to represent the absence of a value.
func None[T any]() Option[T] {
	var value T
	return Option[T]{has: false, value: value}
}
