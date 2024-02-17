package gopts

// Option is a generic type that can be used to represent a value
// that may or may not be present.
type Option[T any] struct {
	value *T
}

// Some returns an Option[T] with a value.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.IsSome()) // true
func Some[T any](value T) Option[T] {
	return Option[T]{&value}
}

// None returns an Option[T] with no value.
// Example:
//
//	opt := None()
//	fmt.Println(opt.IsNone()) // true
func None[T any]() Option[T] {
	return Option[T]{nil}
}

// IsSome returns true if the Options[T] has a value.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.IsSome()) // true
func (o *Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone returns true if the Options[T] has no value.
// Example:
//
//	opt := None()
//	fmt.Println(opt.IsNone()) // true
func (o *Option[T]) IsNone() bool {
	return o.value == nil
}

// Unwrap returns the value of the Options[T].
// If the Options[T] has no value, it panics.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.Unwrap()) // 42
func (o *Option[T]) Unwrap() T {
	if o.value == nil {
		panic("Unwrap a None value")
	}
	return *o.value
}
