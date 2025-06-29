# maybe - Go Optional, Nullable Types and Functional Utilities

`maybe` is a Go package providing type-safe optional values and functional
programming utilities using Go generics. It helps eliminate nil pointer panics
and provides a more expressive way to handle optional values in Go. The package
brings modern constructs like `Option[T]` and `Nullable[T]` to Go with strong
database and JSON integration, along with functional utilities for cleaner data
transformation.

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
