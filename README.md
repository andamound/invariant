# invariant

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library providing type-safe abstractions for more robust programming. The library introduces "invariant" types that enforce guarantees at compile-time and runtime, inspired by Rust's safety patterns.

## Features

- `SafePointer<T>`: A guaranteed non-nil pointer wrapper that ensures pointer validity at creation time
- `Result<T, E>`: Type for error handling similar to Rust's Result, with Ok and Err variants
- Simple, idiomatic API with type parameters (requires Go 1.18+)
- Zero dependencies
- Comprehensive test coverage
- Designed for use in production code where nil pointer panics must be avoided

## Installation

```bash
go get github.com/andamound/invariant
```

## Why Use Invariant?

### The Problem With Go's Approach

In Go, nil pointers are a common source of panics and unexpected behavior:

```go
var user *User = nil
fmt.Println(user.Name) // Panics: nil pointer dereference
```

Function parameters don't indicate whether they can handle nil values:

```go
func processUser(user *User) { 
    // Will this function check for nil? You have to read the implementation
    // or documentation to know for sure.
}
```

Error handling is verbose and can lead to repetitive boilerplate:

```go
value, err := someFunc()
if err != nil {
    // Error handling, potentially repeated many times
    return nil, err
}
// Use value
```

### The Invariant Solution

The `invariant` library addresses these problems with:

1. **Type-safety**: Types that enforce guarantees at the type level
2. **Self-documenting APIs**: Function signatures that clearly indicate constraints
3. **Functional patterns**: More expressive, less error-prone code

## Usage

### SafePointer

SafePointer guarantees that a pointer is never nil:

```go
import "github.com/andamound/invariant"

// Create a SafePointer from a regular pointer
str := "Hello"
safeStr := invariant.SP(&str) // SP is a shorthand for NewSafePointer

// Access the underlying value via the pointer
*safeStr.Some() += " World" 

// Get a copy of the value
value := safeStr.Get() // "Hello World"

// This would panic - SafePointer protects against nil pointers
// var nilPtr *string = nil
// safeStr = invariant.SP(nilPtr) // Panics!
```

### Result

Result provides a more functional approach to error handling:

```go
import (
    "errors"
    "github.com/andamound/invariant"
)

// Return a Result from a function
func divide(a, b int) invariant.Result[float64, error] {
    if b == 0 {
        return invariant.Err[float64, error](errors.New("division by zero"))
    }
    return invariant.Ok[float64, error](float64(a) / float64(b))
}

// Use the Result
result := divide(10, 2)

// Method 1: Check if it's Ok or Err
if result.IsOk() {
    value := result.Unwrap()
    // use value
}

// Method 2: Use Match for more functional style
result.Match(
    func(value float64) { fmt.Println("Result:", value) },
    func(err error) { fmt.Println("Error:", err) },
)

// Method 3: Provide a default value
value := result.UnwrapOr(0.0)
```

## Examples

Here are more examples of how to use the library:

### Using SafePointer with structs

```go
type User struct {
    Name  string
    Email string
}

// Create a function that accepts only non-nil pointers
func updateUser(sp invariant.SafePointer[User]) {
    // Access and modify fields through the pointer
    user := sp.Some()
    user.Name = "Updated Name"
    user.Email = "updated@example.com"
}

// Usage
user := User{Name: "John", Email: "john@example.com"}
updateUser(invariant.SP(&user))
fmt.Println(user) // User will have updated fields
```

### Chaining operations with Result

```go
// Parsing and transforming with Result
func parseAndDouble(input string) invariant.Result[int, error] {
    // Parse the input
    num, err := strconv.Atoi(input)
    if err != nil {
        return invariant.Err[int, error](err)
    }
    
    // Return successful result
    return invariant.Ok[int, error](num)
}

// Usage
result := parseAndDouble("5")
    .Map(func(n int) int {
        return n * 2 // Double the number if parsing was successful
    })

// Handle different outcomes
result.Match(
    func(n int) { fmt.Printf("Result: %d\n", n) }, // Prints "Result: 10"
    func(err error) { fmt.Printf("Error: %v\n", err) },
)
```

### Error handling patterns with Result

```go
// Chain multiple operations that might fail
func processData(data string) invariant.Result[string, error] {
    // First step: validate
    if len(data) == 0 {
        return invariant.Err[string, error](errors.New("empty data"))
    }
    
    // Second step: transform (could fail)
    if strings.Contains(data, "invalid") {
        return invariant.Err[string, error](errors.New("invalid data"))
    }
    
    // Success path
    return invariant.Ok[string, error]("Processed: " + data)
}

// Usage example with error mapping
result := processData(userInput).
    MapErr(func(err error) error {
        return fmt.Errorf("data processing error: %w", err)
    })

// Safely use the result with a default value
output := result.UnwrapOr("default output")
```

## License

[MIT](LICENSE)
