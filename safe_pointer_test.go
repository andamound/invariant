package invariant_test

import (
	"testing"

	"github.com/andamound/invariant"
)

func TestSafePointer(t *testing.T) {
	// Test with a string
	t.Run("String", func(t *testing.T) {
		original := "test"
		sp := invariant.NewSafePointer(&original)

		// Check that Get() returns the value
		if sp.Get() != original {
			t.Errorf("Expected %v, got %v", original, sp.Get())
		}

		// Check that Some() returns the pointer to the value
		if *sp.Some() != original {
			t.Errorf("Expected %v, got %v", original, *sp.Some())
		}

		// Modify via the pointer
		*sp.Some() = "modified"
		if original != "modified" {
			t.Errorf("Expected original value to be modified to %v, got %v", "modified", original)
		}

		// Get should reflect the modified value
		if sp.Get() != "modified" {
			t.Errorf("Expected Get() to return %v, got %v", "modified", sp.Get())
		}
	})

	// Test with a struct
	t.Run("Struct", func(t *testing.T) {
		type testStruct struct {
			Value int
		}

		original := testStruct{Value: 42}
		sp := invariant.SP(&original)

		// Check that Get() returns the value
		if sp.Get().Value != original.Value {
			t.Errorf("Expected %v, got %v", original.Value, sp.Get().Value)
		}

		// Check that Some() returns the pointer to the value
		if sp.Some().Value != original.Value {
			t.Errorf("Expected %v, got %v", original.Value, sp.Some().Value)
		}

		// Modify via the pointer
		sp.Some().Value = 100
		if original.Value != 100 {
			t.Errorf("Expected original.Value to be modified to %v, got %v", 100, original.Value)
		}
	})

	// Test nil pointer handling
	t.Run("Nil pointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic on nil pointer, but no panic occurred")
			}
		}()

		var ptr *int = nil
		_ = invariant.NewSafePointer(ptr) // This should panic
	})
}
