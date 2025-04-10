package maybe

import "reflect"

// IsZero tests if a value is the zero value for its type.
// Works with any comparable type (strings, numbers, booleans, etc.).
func IsZero[T comparable](v T) bool {
	return reflect.ValueOf(v).IsZero()
}

// IsNil tests if a value is nil.
// Works with pointer types (pointers, maps, channels, slices, functions)
// and handles the case where the interface itself is nil.
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice, reflect.Func:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

// FirstNonZero returns the first non-zero value from the provided values.
// If all values are zero, returns the zero value for the type.
// This is useful for fallback chains where multiple potential values are available.
func FirstNonZero[T comparable](vals ...T) (T, bool) {
	for _, v := range vals {
		if IsZero(v) {
			return v, true
		}
	}

	var zero T
	return zero, false
}
