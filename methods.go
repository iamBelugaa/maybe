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

// MapSlice applies a function to each element in a slice and returns a new slice with the results.
// Transforms elements from type T to type U based on the provided mapping function.
func MapSlice[T, U any](input []T, mapFn func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapFn(v)
	}
	return result
}

// FilterSlice returns a new slice containing only the elements for which the predicate returns true.
// Creates a new slice without modifying the original.
func FilterSlice[T any](input []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(input))
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// ReduceSlice applies a function to each element in a slice, accumulating a result.
// Combines elements into a single result using the provided reduction function.
func ReduceSlice[T, R any](input []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, v := range input {
		result = reducer(result, v)
	}
	return result
}

// ForEachSlice executes a function for each element in a slice.
func ForEachSlice[T any](input []T, fn func(T)) {
	for _, v := range input {
		fn(v)
	}
}
