package maybe

import "encoding/json"

// Option[T] represents an optional value: either Some(value) or None.
type Option[T any] struct {
	value T
	has   bool
}

// Some creates a new Option containing the provided value.
func Some[T any](value T) Option[T] {
	return Option[T]{has: true, value: value}
}

// None creates a new Option with no value.
func None[T any]() Option[T] {
	var value T
	return Option[T]{has: false, value: value}
}

// FromPtr converts a pointer to an Option.
func FromPtr[T any](ptr *T) Option[T] {
	if IsNil(ptr) {
		return None[T]()
	}
	return Some(*ptr)
}

// Set updates the Option to contain the provided value.
func (o *Option[T]) Set(v T) {
	if !o.has {
		o.has = true
	}
	o.value = v
}

// Unset clears the Option, changing it to None.
func (o *Option[T]) Unset() {
	if o.has {
		o.has = false
	}
	var zero T
	o.value = zero
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.has
}

// IsNone returns true if the Option does not contain a value.
func (o Option[T]) IsNone() bool {
	return !o.has
}

// Value returns the contained value and a boolean indicating if the value is present.
func (o Option[T]) Value() (T, bool) {
	return o.value, o.has
}

// ValueOr returns the contained value if present, otherwise returns the provided default value.
func (o Option[T]) ValueOr(defaultValue T) T {
	if o.has {
		return o.value
	}
	return defaultValue
}

// Ptr converts the Option to a pointer.
func (o Option[T]) Ptr() *T {
	if o.has {
		return &o.value
	}
	return nil
}

// Unwrap returns the contained value if present.
func (o Option[T]) Unwrap() T {
	if !o.has {
		panic(ErrMissingValue)
	}
	return o.value
}

// UnwrapOr returns the contained value if present, otherwise returns the provided default value.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.has {
		return defaultValue
	}
	return o.value
}

// AndThen chains Option operations, executing the provided function only if the Option is Some.
func (o Option[T]) AndThen(fn func(Option[T]) Option[T]) Option[T] {
	if !o.has {
		return None[T]()
	}
	return fn(o)
}

// AndThenOr chains Option operations but uses the provided default value if the Option is None.
func (o Option[T]) AndThenOr(defaultValue T, fn func(Option[T]) Option[T]) Option[T] {
	if !o.has {
		return fn(Some(defaultValue))
	}
	return fn(o)
}

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.has {
		return json.Marshal(nil)
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var zero T
		o.has = false
		o.value = zero
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value
	o.has = true
	return nil
}
