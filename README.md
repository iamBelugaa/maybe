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
