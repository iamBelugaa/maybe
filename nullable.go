package maybe

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"
)

// Nullable[T] represents a value that might be null.
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
func (n Nullable[T]) ToOption() Option[T] {
	if !n.valid {
		return None[T]()
	}
	return Some(n.value)
}

// Equals compares two Nullable values for equality.
func (n Nullable[T]) Equals(other Nullable[T]) bool {
	if n.valid != other.valid {
		return false
	}

	if !n.valid {
		return true
	}

	// Check if T is comparable.
	if tType := reflect.TypeOf(n.value); !tType.Comparable() {
		return false
	}

	return reflect.DeepEqual(n.value, other.value)
}

// MarshalJSON implements the json.Marshaler interface.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.valid = false
		var zero T
		n.value = zero
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.value = v
	n.valid = true
	return nil
}

// Value implements driver.Valuer to convert Go types to database-compatible values.
func (n Nullable[T]) Value() (driver.Value, error) {
	// Handle NULL case first: if not valid, return nil to represent SQL NULL
	if !n.valid {
		return nil, nil
	}

	// First check if the value implements driver.Valuer itself
	if valuer, ok := any(n.value).(driver.Valuer); ok {
		return valuer.Value()
	}

	// Fast path for common types that don't need conversion
	switch v := any(n.value).(type) {
	case int64, float64, bool, []byte, string, time.Time:
		return v, nil
	}

	rv := reflect.ValueOf(n.value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		{
			return rv.Int(), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			unsignedVal := rv.Uint()
			if unsignedVal > math.MaxInt64 {
				return nil, fmt.Errorf("unsigned integer overflow: %v exceeds int64 maximum", unsignedVal)
			}
			return int64(unsignedVal), nil
		}
	case reflect.Float32:
		{
			return float64(rv.Float()), nil
		}
	}

	return nil, fmt.Errorf("unsupported database type: %T", n.value)
}
