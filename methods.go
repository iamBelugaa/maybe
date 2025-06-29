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
