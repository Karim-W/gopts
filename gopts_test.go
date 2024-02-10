package gopts

import "testing"

func TestSome(t *testing.T) {
	opt := Some(42)
	if !opt.IsSome() {
		t.Error("opt.IsSome() should be true")
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()
	if !opt.IsNone() {
		t.Error("opt.IsNone() should be true")
	}
}

func TestUnwrap(t *testing.T) {
	opt := Some(42)
	if opt.Unwrap() != 42 {
		t.Error("opt.Unwrap() should be 42")
	}
}

func TestUnwrapNone(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("opt.Unwrap() should panic")
		}
	}()
	opt := None[int]()
	opt.Unwrap()
}
