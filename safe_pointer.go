// Package invariant provides types for handling guaranteed non-nil values
package invariant

// SafePointer represents a guaranteed non-nil pointer value.
// It ensures the pointer exists at the time of creation and provides safe access
// to the underlying value.
type SafePointer[T any] struct {
	ptr *T
}

// NewSafePointer creates a new SafePointer instance with the provided pointer.
// This is the recommended way to create a SafePointer.
// Panics if the provided pointer is nil.
func NewSafePointer[T any](ptr *T) SafePointer[T] {
	if ptr == nil {
		panic("invariant.NewSafePointer: nil pointer provided")
	}
	return SafePointer[T]{
		ptr: ptr,
	}
}

// SP is a shorthand function for creating a new SafePointer.
// It's an alias for NewSafePointer for more concise code.
func SP[T any](ptr *T) SafePointer[T] {
	return NewSafePointer(ptr)
}

// Get returns a copy of the value pointed to by the internal pointer.
// This guarantees that the value exists and can be safely accessed.
func (sp SafePointer[T]) Get() T {
	return *sp.ptr
}

// Some returns the underlying pointer, guaranteeing it exists.
// This pointer can be used to access the value safely.
func (sp SafePointer[T]) Some() *T {
	return sp.ptr
}
