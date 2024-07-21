package gopts

import (
	"encoding/json"
	"fmt"
)

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.some {
		return json.Marshal(nil)
	}

	switch v := any(o.value).(type) {
	case json.Marshaler:
		return v.MarshalJSON()
	default:
		return json.Marshal(v)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return fmt.Errorf("gopts: UnmarshalJSON on nil pointer")
	}

	var v T

	if string(data) == "null" {
		(*o).value = v
		(*o).some = false
		return nil
	}

	switch t := any(o.value).(type) {
	case json.Unmarshaler:
		if err := json.Unmarshal(data, t); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
	}

	(*o).value = v
	(*o).some = true

	return nil
}
