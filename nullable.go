package maybe

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

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

// MarshalJSON implements the json.Marshaler interface.
// An invalid Nullable will be marshaled as null.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return nil, ErrMissingValue
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// A null JSON value will be unmarshaled as an invalid Nullable.
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

// Value implements the driver.Valuer interface.
// This method allows Nullable[T] to be used seamlessly with the database/sql package
// when inserting values into a SQL database.
//
// The returned driver.Value must be one of the types supported by database drivers,
func (n Nullable[T]) Value() (driver.Value, error) {
	// If the Nullable is invalid (null), return nil, indicating a SQL NULL value.
	if !n.valid {
		return nil, nil
	}

	// Handle common native types directly — these are already compatible with driver.Value.
	switch v := any(n.Value).(type) {
	case int64, float64, bool, []byte, string, time.Time:
		return v, nil
	}

	// For other types, convert to a driver.Value type
	rv := reflect.ValueOf(n.Value)

	switch rv.Kind() {
	// Convert all integer types to int64, which is universally accepted by drivers.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return int64(rv.Int()), nil
		// Convert unsigned integers to int64 (may truncate large uint64 values).
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), nil
		// Convert float32 to float64, as float64 is accepted by drivers.
	case reflect.Float32:
		return float64(rv.Float()), nil
	}

	return nil, fmt.Errorf("unsupported type %T for database/sql", n.Value)
}
