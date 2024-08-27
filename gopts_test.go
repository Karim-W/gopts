package gopts_test

import (
	"encoding/json"
	"testing"
	"time"

	go_test "github.com/karim-w/go-test"
	"github.com/karim-w/gopts"
)

func TestSome(t *testing.T) {
	opt := gopts.Some(42)
	if !opt.IsSome() {
		t.Error("opt.IsSome() should be true")
	}
}

func TestNone(t *testing.T) {
	opt := gopts.None[int]()
	if !opt.IsNone() {
		t.Error("opt.IsNone() should be true")
	}
}

func TestUnwrap(t *testing.T) {
	opt := gopts.Some(42)
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
	opt := gopts.None[int]()
	opt.Unwrap()
}

func TestGetOrElse(t *testing.T) {
	opt := gopts.Some(42)
	if opt.GetOrElse(0) != 42 {
		t.Error("opt.GetOrElse(0) should be 42")
	}
}

func TestGetOrElseNone(t *testing.T) {
	opt := gopts.None[int]()
	if opt.GetOrElse(0) != 0 {
		t.Error("opt.GetOrElse(0) should be 0")
	}
}

func TestGet(t *testing.T) {
	opt := gopts.Some(42)
	if v, ok := opt.Get(); v != 42 || !ok {
		t.Error("opt.Get() should be 42, true")
	}
}

func TestGetNone(t *testing.T) {
	opt := gopts.None[int]()
	if v, ok := opt.Get(); v != 0 || ok {
		t.Error("opt.Get() should be 0, false")
	}
}

func TestJSONMarshal(t *testing.T) {
	type testStruct struct {
		Opt gopts.Option[int] `json:"opt"`
	}

	opt := gopts.Some(42)
	ts := testStruct{Opt: opt}
	data, err := json.Marshal(ts)
	if err != nil {
		t.Error(err)
	}

	expected := `{"opt":42}`
	if string(data) != expected {
		t.Errorf("got %s, expected %s", data, expected)
	}

	opt = gopts.None[int]()
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
		Opt gopts.Option[int] `json:"opt"`
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
		Opt gopts.Option[int] `json:"opt"`
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

func TestNullableSQLScan(t *testing.T) {
	db, cleanup := go_test.InitDockerPostgresSQLDBTest(t)
	defer cleanup()

	var opt gopts.Option[int]

	const migration = `
	CREATE TABLE test (
		id SERIAL PRIMARY KEY,
		opt INT NULL
	);`

	_, err := db.Exec(migration)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test (opt) VALUES ($1)", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = db.QueryRow("SELECT opt FROM test").Scan(&opt)
	if err != nil {
		t.Fatal(err)
	}

	if opt.IsSome() {
		t.Fatal("opt should be None")
	}

	_, err = db.Exec("INSERT INTO test (opt) VALUES ($1)", 42)
	if err != nil {
		t.Fatal(err)
	}

	err = db.QueryRow("SELECT opt FROM test WHERE opt IS NOT NULL").Scan(&opt)
	if err != nil {
		t.Fatal(err)
	}

	if opt.IsNone() {
		t.Fatal("opt should be Some(42)")
	}

	if v, ok := opt.Get(); v != 42 || !ok {
		t.Errorf("got %d, expected 42", v)
	}
}

func TestNullableSQLAllTypesScan(t *testing.T) {
	db, cleanup := go_test.InitDockerPostgresSQLDBTest(t)
	defer cleanup()

	const migration = `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		age INT NULL,
		salary BIGINT NULL,
		name TEXT NULL,
		active BOOLEAN NULL,
		weight REAL NULL,
		created_at TIMESTAMP WITH TIME ZONE NULL
	);`

	_, err := db.Exec(migration)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		"INSERT INTO users (age, salary, name, active, weight,created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	var optAge gopts.Option[int]
	var optSalary gopts.Option[int64]
	var optWeight gopts.Option[float64]
	var optName gopts.Option[string]
	var optActive gopts.Option[bool]
	var optTime gopts.Option[time.Time]

	// Example database query
	row := db.QueryRow(
		"SELECT age, salary, name, active, weight ,created_at FROM users WHERE id = $1",
		1,
	)

	err = row.Scan(&optAge, &optWeight, &optName, &optActive, &optWeight, &optTime)
	if err != nil {
		t.Fatal(err)
		return
	}

	if optAge.IsSome() {
		t.Fatal("optAge should be None")
	}

	if optSalary.IsSome() {
		t.Fatal("optSalary should be None")
	}

	if optWeight.IsSome() {
		t.Fatal("optWeight should be None")
	}

	if optName.IsSome() {
		t.Fatal("optName should be None")
	}

	if optActive.IsSome() {
		t.Fatal("optActive should be None")
	}

	if optTime.IsSome() {
		t.Fatal("optTime should be None")
	}

	_, err = db.Exec(
		"INSERT INTO users (age, salary, name, active,weight ,created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		42,
		42,
		"foo",
		true,
		42.4,
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	row = db.QueryRow(
		"SELECT age, salary, name, active, weight,created_at FROM users WHERE id = $1",
		2,
	)

	err = row.Scan(&optAge, &optSalary, &optName, &optActive, &optWeight, &optTime)
	if err != nil {
		t.Fatal(err)
		return
	}

	if optAge.IsNone() {
		t.Fatal("optAge should be Some(42)")
	}

	if optSalary.IsNone() {
		t.Fatal("optSalary should be Some(42)")
	}

	if optWeight.IsNone() {
		t.Fatal("optWeight should be Some(42.4)")
	}

	if optName.IsNone() {
		t.Fatal("optName should be Some(\"foo\")")
	}

	if optActive.IsNone() {
		t.Fatal("optActive should be Some(true)")
	}

	if optTime.IsNone() {
		t.Fatal("optTime should be Some(time.Now())")
		if optTime.Unwrap().IsZero() {
			t.Fatal("optTime should not be zero")
		}
	}
}

func TestNullableSQLAllTypesInsert(t *testing.T) {
	db, cleanup := go_test.InitDockerPostgresSQLDBTest(t)
	defer cleanup()

	const migration = `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		age INT NULL,
		salary BIGINT NULL,
		name TEXT NULL,
		active BOOLEAN NULL,
		weight REAL NULL,
		created_at TIMESTAMP WITH TIME ZONE NULL
	);`

	_, err := db.Exec(migration)
	if err != nil {
		t.Fatal(err)
	}

	var optAge gopts.Option[int]
	var optSalary gopts.Option[int64]
	var optWeight gopts.Option[float64]
	var optName gopts.Option[string]
	var optActive gopts.Option[bool]
	var optTime gopts.Option[time.Time]
	_, err = db.Exec(
		"INSERT INTO users (age, salary, name, active, weight,created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		optAge,
		optSalary,
		optName,
		optActive,
		optWeight,
		optTime,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Example database query
	row := db.QueryRow(
		"SELECT age, salary, name, active, weight ,created_at FROM users WHERE id = $1",
		1,
	)

	err = row.Scan(&optAge, &optWeight, &optName, &optActive, &optWeight, &optTime)
	if err != nil {
		t.Fatal(err)
		return
	}

	if optAge.IsSome() {
		t.Fatal("optAge should be None")
	}

	if optSalary.IsSome() {
		t.Fatal("optSalary should be None")
	}

	if optWeight.IsSome() {
		t.Fatal("optWeight should be None")
	}

	if optName.IsSome() {
		t.Fatal("optName should be None")
	}

	if optActive.IsSome() {
		t.Fatal("optActive should be None")
	}

	if optTime.IsSome() {
		t.Fatal("optTime should be None")
	}

	_, err = db.Exec(
		"INSERT INTO users (age, salary, name, active,weight ,created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		gopts.Some(42),
		gopts.Some(42),
		gopts.Some("foo"),
		gopts.Some(true),
		gopts.Some(42.4),
		gopts.Some(time.Now()),
	)
	if err != nil {
		t.Fatal(err)
	}

	row = db.QueryRow(
		"SELECT age, salary, name, active, weight,created_at FROM users WHERE id = $1",
		2,
	)

	err = row.Scan(&optAge, &optSalary, &optName, &optActive, &optWeight, &optTime)
	if err != nil {
		t.Fatal(err)
		return
	}

	if optAge.IsNone() {
		t.Fatal("optAge should be Some(42)")
	}

	if optSalary.IsNone() {
		t.Fatal("optSalary should be Some(42)")
	}

	if optWeight.IsNone() {
		t.Fatal("optWeight should be Some(42.4)")
	}

	if optName.IsNone() {
		t.Fatal("optName should be Some(\"foo\")")
	}

	if optActive.IsNone() {
		t.Fatal("optActive should be Some(true)")
	}

	if optTime.IsNone() {
		t.Fatal("optTime should be Some(time.Now())")
		if optTime.Unwrap().IsZero() {
			t.Fatal("optTime should not be zero")
		}
	}
}
