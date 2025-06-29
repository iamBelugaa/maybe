# maybe - A Go Package for Optionals, Nullables, and Nil-Safe Data Handling

`maybe` is a Go package providing type-safe optional values and functional
programming utilities using Go generics. It helps eliminate nil pointer panics
and provides a more expressive way to handle optional values in Go.

## Features

- `Option[T]`: Type-safe optional values (Some/None) for any type.
- `Nullable[T]`: Optional values specifically designed for database and JSON
  null values.
- Functional programming utilities (Map, Filter, Reduce, etc.).
- JSON marshaling/unmarshaling support.
- SQL database integration with `database/sql` compatibility.
- Zero dependencies beyond the standard library.
- Fully compatible with Go generics (Go 1.18+).

## Why Use `maybe`?

- **Type Safety**: Avoid nil pointer panics and make optional values explicit.
- **Expressiveness**: Clearer intent than nil pointers or pointer returns.
- **Composability**: Functional utilities for working with collections.
- **Interoperability**: Easy conversion between options, nullables, and
  pointers.
- **Maintainability**: More readable code with explicit handling of missing
  values.
- **Database Integration**: First-class support for handling SQL NULL values.
- **JSON Compatibility**: Seamless handling of missing or null JSON fields.

## Installation

```bash
go get github.com/iamNilotpal/maybe
```

## Core Types

### Option

`Option[T]` represents an optional value: either `Some(value)` or `None`. It
provides a type-safe alternative to nil pointers and helps avoid nil pointer
panics. The zero value of `Option` is `None` (no value present).

Key operations:

- `Some(value)`: Create an Option with a value.
- `None[T]()`: Create an Option with no value.
- `FromPtr(ptr)`: Convert a pointer to an Option.
- `IsSome()`: Check if value is present.
- `IsNone()`: Check if value is absent.
- `Value()`: Get the value and a success boolean.
- `ValueOr(default)`: Get the value or a default.
- `Unwrap()`: Get the value or panic.
- `Ptr()`: Convert to a pointer (nil if None).

### Nullable

`Nullable[T]` represents a value that might be null, designed specifically for
handling null values in databases and JSON. Unlike `Option[T]`, which is for
general-purpose optional values, `Nullable[T]` is optimized for scenarios
involving external systems that use null values.

Key operations:

- `NullableOf(value)`: Create a valid Nullable.
- `Null[T]()`: Create a null Nullable.
- `NullableFromPtr(ptr)`: Create a Nullable from a pointer.
- `IsNull()`: Check if null.
- `IsValid()`: Check if not null.
- `Extract()`: Get the value and validity.
- `ExtractOr(default)`: Get the value or a default.
- `ToPtr()`: Convert to pointer (nil if null).
- `ToOption()`: Convert to Option type.

## Usage Examples

### Basic Operations

```go
package main

import (
	"fmt"

	"github.com/iamNilotpal/maybe"
)

func main() {
	// Working with Option
	name := maybe.Some("Nilotpal")
	emptyName := maybe.None[string]()

	// Checking value presence
	fmt.Println("Has name:", name.IsSome())            // true
	fmt.Println("Has empty name:", emptyName.IsNone()) // true

	// Safe access patterns
	if value, ok := name.Value(); ok {
		fmt.Println("Name:", value) // "Nilotpal"
	}

	// Default values
	fmt.Println("Name or default:", name.ValueOr("Anonymous"))            // "Nilotpal"
	fmt.Println("Empty name or default:", emptyName.ValueOr("Anonymous")) // "Anonymous"

	// Unwrap (safe only when you're certain the value exists)
	fmt.Println("Unwrapped name:", name.Unwrap()) // "Nilotpal"
	// emptyName.Unwrap() would panic with ErrMissingValue

	// Get or set a value
	emptyName.Set("Bob")
	fmt.Println("Name after set:", emptyName.ValueOr("")) // "Bob"

	emptyName.Unset()
	fmt.Println("Is name none after unset:", emptyName.IsNone()) // true

	// Convert to/from pointers
	var _ *string = name.Ptr()           // Pointer to "Nilotpal"
	var nilPtr *string = emptyName.Ptr() // nil pointer

	someStr := "Hello"
	_ = maybe.FromPtr(&someStr) // Some("Hello")
	_ = maybe.FromPtr(nilPtr)   // None

	// Working with Nullable (for database/JSON)
	userID := maybe.NullableOf(123)
	noID := maybe.Null[int]()

	fmt.Println("Has ID:", userID.IsValid()) // true
	fmt.Println("No ID:", noID.IsNull())     // true

	// Extract value from Nullable
	if val, ok := userID.Extract(); ok {
		fmt.Println("User ID:", val) // 123
	}

	// Default values with Nullable
	fmt.Println("ID or default:", userID.ExtractOr(0))    // 123
	fmt.Println("No ID or default:", noID.ExtractOr(999)) // 999

	// Convert between Option and Nullable
	_ = userID.ToOption() // Some(123)

	// Create Nullable from pointer
	ptrID := userID.ToPtr()          // Pointer to 123
	_ = maybe.NullableFromPtr(ptrID) // Valid Nullable with 123

	// Using with zero values
	zeroInt := maybe.NullableOf(0)                     // Valid Nullable containing 0
	fmt.Println("Is zero int null?", zeroInt.IsNull()) // false

	// Equality check
	anotherZero := maybe.NullableOf(0)
	fmt.Println("Equal zero values:", zeroInt.Equals(anotherZero))   // true
	fmt.Println("Equal to different value:", zeroInt.Equals(userID)) // false
	fmt.Println("Both null equal:", noID.Equals(maybe.Null[int]()))  // true
}
```

### JSON Handling

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/iamNilotpal/maybe"
)

type Person struct {
    Name     string                 `json:"name"`
    Age      maybe.Option[int]      `json:"age,omitempty"`   // Omitted if not present
    Phone    maybe.Nullable[string] `json:"phone"`           // Explicit null if not present
    Address  maybe.Option[Address]  `json:"address,omitempty"`
}

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

func main() {
    // Creating a person with some fields missing
    person := Person{
        Name:  "John Doe",
        Age:   maybe.Some(30),
        Phone: maybe.Null[string](),  // Explicitly null
        // Address is implicitly None
    }

    // Marshal to JSON
    data, _ := json.MarshalIndent(person, "", "  ")
    fmt.Println(string(data))
    // Output:
    // {
    //   "name": "John Doe",
    //   "age": 30,
    //   "phone": null
    // }

    // JSON with explicit null vs missing field
    jsonData := []byte(`{
        "name": "Nilotpal Deka",
        "age": null,
        "phone": "555-1234",
        "address": {
            "street": "123 Main St",
            "city": "Guwahati",
            "country": "IND"
        }
    }`)

    var anotherPerson Person
    json.Unmarshal(jsonData, &anotherPerson)

    // Age was null in JSON, so it's None in our struct
    fmt.Println("Has age:", anotherPerson.Age.IsSome())  // false

    // Phone was present, so it's a valid Nullable
    fmt.Println("Has phone:", anotherPerson.Phone.IsValid())  // true
    if phone, ok := anotherPerson.Phone.Extract(); ok {
        fmt.Println("Phone:", phone)  // "555-1234"
    }

    // Address was present as an object
    if address, ok := anotherPerson.Address.Value(); ok {
        fmt.Println("City:", address.City)  // "Guwahati"
    }

    // Error handling with MarshalJSON
    badJSON := []byte(`{"name": "Bad Data", "phone": ["invalid"]}`)
    var badPerson Person
    err := json.Unmarshal(badJSON, &badPerson)
    if err != nil {
        fmt.Println("JSON error:", err)
    }
}
```

## API Reference

### Option Methods

- **`Some[T](value T) Option[T]`**: Creates a new Option containing the provided
  value.
- **`None[T]() Option[T]`**: Creates a new Option with no value.
- **`FromPtr[T](ptr *T) Option[T]`**: Converts a pointer to an Option. Returns
  `Some(value)` if pointer is non-nil, or `None` if pointer is nil.
- **`Set(v T)`**: Updates the Option to contain the provided value. Changes None
  to Some(value) or updates an existing Some value.
- **`Unset()`**: Clears the Option, changing it to None. The contained value is
  set to the zero value.
- **`IsSome() bool`**: Returns true if the Option contains a value.
- **`IsNone() bool`**: Returns true if the Option does not contain a value.
- **`Value() (T, bool)`**: Returns the contained value and a boolean indicating
  if the value is present.
- **`ValueOr(defaultValue T) T`**: Returns the contained value if present,
  otherwise returns the provided default value.
- **`Ptr() *T`**: Converts the Option to a pointer. Returns a pointer to the
  value if Some, or nil if None.
- **`Unwrap() T`**: Returns the contained value if present. Panics with
  ErrMissingValue if the Option is None.
- **`UnwrapOr(defaultValue T) T`**: Returns the contained value if present,
  otherwise returns the provided default value.
- **`AndThen(fn func(Option[T]) Option[T]) Option[T]`**: Chains Option
  operations, executing the provided function only if the Option is Some.
- **`AndThenOr(defaultValue T, fn func(Option[T]) Option[T]) Option[T]`**:
  Chains Option operations but uses the provided default value if the Option is
  None.
- **`MarshalJSON() ([]byte, error)`**: Marshals the Option to JSON. Returns an
  error if the Option is None.
- **`UnmarshalJSON(data []byte) error`**: Unmarshal JSON data into the Option,
  setting it to None if the JSON value is null.

### Nullable Methods

- **`NullableOf[T](value T) Nullable[T]`**: Creates a valid Nullable with the
  provided value.
- **`Null[T]() Nullable[T]`**: Creates an invalid (null) Nullable.
- **`NullableFromPtr[T](ptr *T) Nullable[T]`**: Creates a Nullable from a
  pointer. If the pointer is nil, returns an invalid Nullable.
- **`IsNull() bool`**: Returns true if this represents a null value.
- **`IsValid() bool`**: Returns true if this represents a non-null value.
- **`Extract() (T, bool)`**: Returns the contained value and a boolean
  indicating if the value is valid.
- **`ExtractOr(defaultVal T) T`**: Returns the value if valid, otherwise returns
  the default.
- **`ToPtr() *T`**: Converts to a pointer, which will be nil if the value is
  null.
- **`ToOption() Option[T]`**: Converts Nullable to an Option type.
- **`Equals(other Nullable[T]) bool`**: Compares two Nullable values for
  equality.
- **`MarshalJSON() ([]byte, error)`**: Implements the json.Marshaler interface.
  An invalid Nullable will be marshaled as null.
- **`UnmarshalJSON(data []byte) error`**: Implements the json.Unmarshaler
  interface. A null JSON value will be unmarshaled as an invalid Nullable.

### Utility Functions Reference

- **`IsZero[T comparable](v T) bool`**: Tests if a value is the zero value for
  its type.
- **`IsNil(i any) bool`**: Tests if a value is nil. Works with pointer types and
  handles the case where the interface itself is nil.
- **`FirstNonZero[T comparable](vals ...T) (T, bool)`**: Returns the first
  non-zero value from the provided values.
- **`MapSlice[T, U any](input []T, mapFn func(T) U) []U`**: Applies a function
  to each element in a slice and returns a new slice with the results.
- **`FilterSlice[T any](input []T, predicate func(T) bool) []T`**: Returns a new
  slice containing only the elements for which the predicate returns true.
- **`ReduceSlice[T, R any](input []T, initial R, reducer func(R, T) R) R`**:
  Applies a function to each element in a slice, accumulating a result.
- **`ForEachSlice[T any](input []T, fn func(T))`**: Executes a function for each
  element in a slice.
- **`CollectOptions[T any](options []Option[T]) Option[[]T]`**: Transforms a
  slice of Options into an Option containing a slice of all Some values.
- **`FilterSomeOptions[T any](options []Option[T]) []T`**: Returns a slice
  containing only the values from non-empty Options.
- **`PartitionOptions[T any](options []Option[T]) (values []T, noneIndices []int)`**:
  Separates a slice of Options into values from Some options and indices of None
  options.
- **`TryMap[T, U any](input []T, fn func(T) Option[U]) Option[[]U]`**: Applies a
  function that might fail to each element in a slice.

## License

MIT License - see [LICENSE](LICENSE) for details.
