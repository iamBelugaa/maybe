package maybe

import "reflect"

// Nullable[T] represents a value that might be null.
// Unlike Option, Nullable is specifically designed for handling
// null values in external systems like databases and JSON APIs.
// This is particularly useful for database operations where NULL
// values are common and need distinct handling.
type Nullable[T any] struct {
	value T
	valid bool
}

// NullableOf creates a valid Nullable with the provided value.
func NullableOf[T any](value T) Nullable[T] {
	return Nullable[T]{value: value, valid: true}
}

// Null creates an invalid (null) Nullable.
func Null[T any]() Nullable[T] {
	var zero T
	return Nullable[T]{value: zero, valid: false}
}

// NullableFromPtr creates a Nullable from a pointer.
// If the pointer is nil, returns an invalid Nullable.
func NullableFromPtr[T any](ptr *T) Nullable[T] {
	if IsNil(ptr) {
		return Null[T]()
	}
	return NullableOf(*ptr)
}

// IsNull returns true if this represents a null value.
func (n Nullable[T]) IsNull() bool {
	return !n.valid
}

// IsValid returns true if this represents a non-null value.
func (n Nullable[T]) IsValid() bool {
	return n.valid
}

// Value returns the contained value and a boolean indicating if the value is valid.
// If the Nullable is null, returns the zero value of T and false.
func (n Nullable[T]) Extract() (T, bool) {
	return n.value, n.valid
}

// ExtractOr returns the value if valid, otherwise returns the default.
func (n Nullable[T]) ExtractOr(defaultVal T) T {
	if n.valid {
		return n.value
	}
	return defaultVal
}

// ToPtr converts to a pointer, which will be nil if the value is null.
func (n Nullable[T]) ToPtr() *T {
	if !n.valid {
		return nil
	}
	return &n.value
}

// ToOption converts Nullable to an Option type.
// This allows for interoperability between the two optional value representations.
func (n Nullable[T]) ToOption() Option[T] {
	if !n.valid {
		return None[T]()
	}
	return Some(n.value)
}

// Equals compares two Nullable values for equality.
// Two Nullable values are equal if:
//  1. Both are null, or
//  2. Both are valid and contain equal values
func (n Nullable[T]) Equals(other Nullable[T]) bool {
	if n.valid != other.valid {
		return false
	}

	if !n.valid {
		return true // Both are null
	}

	return reflect.ValueOf(n.value).Equal(reflect.ValueOf(other.value))
}
