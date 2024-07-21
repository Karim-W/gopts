package gopts

import (
	"encoding/json"
	"testing"
)

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

func TestJSONMarshal(t *testing.T) {
	type testStruct struct {
		Opt Option[int] `json:"opt"`
	}

	opt := Some(42)
	ts := testStruct{Opt: opt}
	data, err := json.Marshal(ts)
	if err != nil {
		t.Error(err)
	}

	expected := `{"opt":42}`
	if string(data) != expected {
		t.Errorf("got %s, expected %s", data, expected)
	}

	opt = None[int]()
	ts = testStruct{Opt: opt}
	data, err = json.Marshal(ts)
	if err != nil {
		t.Error(err)
	}

	expected = `{"opt":null}`
	if string(data) != expected {
		t.Errorf("got %s, expected %s", data, expected)
	}
}

func TestJSONUnmarshal(t *testing.T) {
	type testStruct struct {
		Opt Option[int] `json:"opt"`
	}

	data := []byte(`{"opt":42}`)
	var ts testStruct
	err := json.Unmarshal(data, &ts)
	if err != nil {
		t.Error(err)
	}

	if v, ok := ts.Opt.Get(); v != 42 || !ok {
		t.Errorf("got %d, expected 42", v)
	}

	var tsNull testStruct

	data = []byte(`{"opt":null}`)
	err = json.Unmarshal(data, &tsNull)
	if err != nil {
		t.Error(err)
	}

	if v, ok := tsNull.Opt.Get(); v != 0 || ok {
		t.Errorf("got %d, expected 0", v)
	}
}

func TestJSONUnmarshalError(t *testing.T) {
	type testStruct struct {
		Opt Option[int] `json:"opt"`
	}

	data := []byte(`{"opt":42}`)
	var ts testStruct
	err := json.Unmarshal(data, &ts)
	if err != nil {
		t.Error(err)
	}

	if v, ok := ts.Opt.Get(); v != 42 || !ok {
		t.Errorf("got %d, expected 42", v)
	}

	data = []byte(`{"opt":"foo"}`)
	err = json.Unmarshal(data, &ts)
	if err == nil {
		t.Error("expected error")
	}
}
