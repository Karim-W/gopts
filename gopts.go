package gopts

// Option is a generic type that can be used to represent a value
// that may or may not be present.
type Option[T any] struct {
	value T
	some  bool
}

// Option is a generic interface that can be used to represent a value
// that may or may not be present.
// type Option[T any] interface {
// 	IsSome() bool
// 	IsNone() bool
// 	Unwrap() T
// 	GetOrElse(T) T
// 	Get() (T, bool)
// }

// Some returns an Option[T] with a value.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.IsSome()) // true
func Some[T any](value T) Option[T] {
	return Option[T]{value, true}
}

// None returns an Option[T] with no value.
// Example:
//
//	opt := None()
//	fmt.Println(opt.IsNone()) // true
func None[T any]() Option[T] {
	return Option[T]{some: false}
}

// IsSome returns true if the Options[T] has a value.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.IsSome()) // true
func (o *Option[T]) IsSome() bool {
	return o.some
}

// IsNone returns true if the Options[T] has no value.
// Example:
//
//	opt := None()
//	fmt.Println(opt.IsNone()) // true
func (o *Option[T]) IsNone() bool {
	return o.some == false
}

// Unwrap returns the value of the Options[T].
// If the Options[T] has no value, it panics.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.Unwrap()) // 42
func (o *Option[T]) Unwrap() (res T) {
	if o.some == false {
		panic("Unwrap called on None")
	}

	return o.value
}

// GetOrElse returns the value of the Options[T].
// If the Options[T] has no value, it returns the default value.
// Example:
//
//	opt := Some(42)
//	fmt.Println(opt.GetOrElse(0)) // 42
func (o *Option[T]) GetOrElse(defaultValue T) (res T) {
	if o.some == false {
		return defaultValue
	}

	return o.value
}

// Get returns the value of the Options[T] and a boolean indicating
// if the value is present.
// Example:
//
//	opt := Some(42)
//	val, ok := opt.Get()
//	fmt.Println(val, ok) // 42, true
func (o *Option[T]) Get() (res T, ok bool) {
	if o.some == false {
		return
	}

	return o.value, true
}
