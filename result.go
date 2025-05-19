// Package invariant provides types for handling guaranteed non-nil values
package invariant

import (
	"fmt"
	"reflect"
)

// Result represents a value that is either successful (Ok) or an error (Err).
// It's similar to Rust's Result type and ensures that exactly one of Ok or Err
// is non-nil at any given time.
type Result[T any, E error] struct {
	ok   *T
	err  E
	isOk bool // True if this is an Ok result, false if it's an Err result
}

// Ok creates a new Result with a successful value.
// The error will be nil.
// Panics if the value is nil (for pointer or interface types).
func Ok[T any, E error](value T) Result[T, E] {
	// For pointer or interface types, check if value is nil
	if isNil(value) {
		panic("invariant.Ok: nil value provided")
	}

	var err E // Zero value for the error type
	return Result[T, E]{
		ok:   &value,
		err:  err,
		isOk: true,
	}
}

// Err creates a new Result with an error.
// The success value will be nil.
// Panics if the error is nil.
func Err[T any, E error](err E) Result[T, E] {
	// Check if err is nil
	if isNilError(err) {
		panic("invariant.Err: nil error provided")
	}

	return Result[T, E]{
		ok:   nil,
		err:  err,
		isOk: false,
	}
}

// isNil is a helper function to check if a value is nil
// for interface and pointer types
func isNil(v any) bool {
	if v == nil {
		return true
	}

	// Use reflection for pointer/interface types
	val := reflect.ValueOf(v)
	kind := val.Kind()
	if (kind == reflect.Ptr || kind == reflect.Interface ||
		kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan ||
		kind == reflect.Func) && val.IsNil() {
		return true
	}

	return false
}

// isNilError is a helper function to check if an error is nil
func isNilError(err error) bool {
	return err == nil || reflect.ValueOf(err).IsNil()
}

// IsOk returns true if the Result contains a success value.
func (r Result[T, E]) IsOk() bool {
	return r.isOk
}

// IsErr returns true if the Result contains an error.
func (r Result[T, E]) IsErr() bool {
	return !r.isOk
}

// Unwrap returns the contained Ok value or panics if the result is an error.
// It's similar to Rust's unwrap() method.
func (r Result[T, E]) Unwrap() T {
	if !r.isOk {
		panic(fmt.Sprintf("invariant.Result.Unwrap: called on Err value: %v", r.err))
	}
	return *r.ok
}

// UnwrapOr returns the contained Ok value or the provided default value if the result is an error.
func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if !r.isOk {
		return defaultValue
	}
	return *r.ok
}

// UnwrapErr returns the contained Err value or panics if the result is a success.
func (r Result[T, E]) UnwrapErr() E {
	if r.isOk {
		panic("invariant.Result.UnwrapErr: called on Ok value")
	}
	return r.err
}

// Map applies a function to the contained value if the result is Ok,
// otherwise returns the Err value unchanged.
func (r Result[T, E]) Map(f func(T) T) Result[T, E] {
	if r.isOk {
		newValue := f(*r.ok)
		return Ok[T, E](newValue)
	}
	return r
}

// MapErr applies a function to the contained error if the result is Err,
// otherwise returns the Ok value unchanged.
func (r Result[T, E]) MapErr(f func(E) E) Result[T, E] {
	if !r.isOk {
		return Err[T, E](f(r.err))
	}
	return r
}

// Match executes okFn if the result is Ok, or errFn if the result is Err.
// This provides a way to handle both cases with a single function call.
func (r Result[T, E]) Match(okFn func(T), errFn func(E)) {
	if r.isOk {
		okFn(*r.ok)
	} else {
		errFn(r.err)
	}
}
