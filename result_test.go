package invariant_test

import (
	"errors"
	"testing"

	"github.com/andamound/invariant"
)

func TestResult(t *testing.T) {
	// Test Ok result
	t.Run("Ok result", func(t *testing.T) {
		result := invariant.Ok[string, error]("success")

		// Check state
		if !result.IsOk() {
			t.Error("Expected IsOk() to be true")
		}
		if result.IsErr() {
			t.Error("Expected IsErr() to be false")
		}

		// Check Unwrap
		if result.Unwrap() != "success" {
			t.Errorf("Expected Unwrap() to return 'success', got '%v'", result.Unwrap())
		}

		// Check UnwrapOr
		if result.UnwrapOr("default") != "success" {
			t.Errorf("Expected UnwrapOr() to return 'success', got '%v'", result.UnwrapOr("default"))
		}

		// Check Map
		mappedResult := result.Map(func(s string) string {
			return s + "!"
		})
		if mappedResult.Unwrap() != "success!" {
			t.Errorf("Expected Map() to return 'success!', got '%v'", mappedResult.Unwrap())
		}

		// Check Match
		var matchResult string
		result.Match(
			func(s string) { matchResult = "Ok: " + s },
			func(err error) { matchResult = "Err: " + err.Error() },
		)
		if matchResult != "Ok: success" {
			t.Errorf("Expected Match() to set matchResult to 'Ok: success', got '%v'", matchResult)
		}
	})

	// Test Err result
	t.Run("Err result", func(t *testing.T) {
		testErr := errors.New("test error")
		result := invariant.Err[string, error](testErr)

		// Check state
		if result.IsOk() {
			t.Error("Expected IsOk() to be false")
		}
		if !result.IsErr() {
			t.Error("Expected IsErr() to be true")
		}

		// Check UnwrapErr
		if result.UnwrapErr().Error() != "test error" {
			t.Errorf("Expected UnwrapErr() to return 'test error', got '%v'", result.UnwrapErr())
		}

		// Check UnwrapOr
		if result.UnwrapOr("default") != "default" {
			t.Errorf("Expected UnwrapOr() to return 'default', got '%v'", result.UnwrapOr("default"))
		}

		// Check MapErr
		mappedResult := result.MapErr(func(err error) error {
			return errors.New("mapped: " + err.Error())
		})
		if mappedResult.UnwrapErr().Error() != "mapped: test error" {
			t.Errorf("Expected MapErr() to return 'mapped: test error', got '%v'", mappedResult.UnwrapErr())
		}

		// Check Match
		var matchResult string
		result.Match(
			func(s string) { matchResult = "Ok: " + s },
			func(err error) { matchResult = "Err: " + err.Error() },
		)
		if matchResult != "Err: test error" {
			t.Errorf("Expected Match() to set matchResult to 'Err: test error', got '%v'", matchResult)
		}
	})

	// Test Unwrap panic on Err
	t.Run("Unwrap panic on Err", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected Unwrap() to panic on Err result")
			}
		}()

		result := invariant.Err[int, error](errors.New("test error"))
		_ = result.Unwrap() // This should panic
	})

	// Test UnwrapErr panic on Ok
	t.Run("UnwrapErr panic on Ok", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected UnwrapErr() to panic on Ok result")
			}
		}()

		result := invariant.Ok[int, error](42)
		_ = result.UnwrapErr() // This should panic
	})
}

func TestResultNilHandling(t *testing.T) {
	// Test nil value in Ok for pointer type
	t.Run("Nil value in Ok for pointer type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when providing nil to Ok for a pointer type")
			}
		}()

		var nilPtr *string = nil
		_ = invariant.Ok[*string, error](nilPtr) // This should panic
	})

	// Test nil error in Err
	t.Run("Nil error in Err", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when providing nil to Err")
			}
		}()

		var nilErr error = nil
		_ = invariant.Err[int, error](nilErr) // This should panic
	})
}
