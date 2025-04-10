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

// FromPtr converts a pointer to an Option.
// Returns Some(value) if the pointer is non-nil, or None if the pointer is nil.
func FromPtr[T any](ptr *T) Option[T] {
	if IsNil(ptr) {
		return None[T]()
	}
	return Some(*ptr)
}

// Set updates the Option to contain the provided value.
// Changes None to Some(value) or updates an existing Some value.
func (o *Option[T]) Set(v T) {
	if !o.has {
		o.has = true
	}
	o.value = v
}

// Unset clears the Option, changing it to None.
// The contained value (if any) is set to the zero value of type T.
func (o *Option[T]) Unset() {
	if o.has {
		o.has = false
	}
	var zero T
	o.value = zero
}
