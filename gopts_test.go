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

func TestGetOrElse(t *testing.T) {
	opt := Some(42)
	if opt.GetOrElse(0) != 42 {
		t.Error("opt.GetOrElse(0) should be 42")
	}
}

func TestGetOrElseNone(t *testing.T) {
	opt := None[int]()
	if opt.GetOrElse(0) != 0 {
		t.Error("opt.GetOrElse(0) should be 0")
	}
}

func TestGet(t *testing.T) {
	opt := Some(42)
	if v, ok := opt.Get(); v != 42 || !ok {
		t.Error("opt.Get() should be 42, true")
	}
}

func TestGetNone(t *testing.T) {
	opt := None[int]()
	if v, ok := opt.Get(); v != 0 || ok {
		t.Error("opt.Get() should be 0, false")
	}
}
