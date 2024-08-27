package gopts

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Scan implements the sql.Scanner interface for Option
func (o *Option[T]) Scan(value interface{}) error {
	if value == nil {
		o.some = false
		return nil
	}

	switch any(o.value).(type) {
	case int:
		if val, ok := value.(int64); ok {
			o.value = any(int(val)).(T)
		} else {
			return errors.New("expected int64 from database for int type")
		}
	case int64:
		if val, ok := value.(int64); ok {
			o.value = any(val).(T)
		} else {
			return errors.New("expected int64 from database for int64 type")
		}
	case float32:
		if val, ok := value.(float64); ok {
			o.value = any(float32(val)).(T)
		} else {
			return errors.New("expected float64 from database for float32 type")
		}
	case float64:
		if val, ok := value.(float64); ok {
			o.value = any(val).(T)
		} else {
			return errors.New("expected float64 from database for float64 type")
		}
	case string:
		if val, ok := value.(string); ok {
			o.value = any(val).(T)
		} else if val, ok := value.([]byte); ok {
			o.value = any(string(val)).(T)
		} else {
			return errors.New("expected string or []byte from database for string type")
		}
	case []byte:
		if val, ok := value.([]byte); ok {
			o.value = any(val).(T)
		} else {
			return errors.New("expected []byte from database for []byte type")
		}
	case bool:
		if val, ok := value.(bool); ok {
			o.value = any(val).(T)
		} else {
			return errors.New("expected bool from database for bool type")
		}
	case time.Time:
		if val, ok := value.(time.Time); ok {
			o.value = any(val).(T)
		} else {
			return errors.New("expected time.Time from database for time.Time type")
		}
	default:
		return errors.New("unsupported scan type")
	}

	o.some = true
	return nil
}

// Value implements the driver.Valuer interface for Option[T]
func (o Option[T]) Value() (driver.Value, error) {
	if !o.some {
		// Return nil to represent a SQL NULL
		return nil, nil
	}

	switch v := any(o.value).(type) {
	case int:
		return int64(v), nil
	case float32:
		return float64(v), nil
	case driver.Valuer:
		return v.Value()
	default:
		return v, nil
	}
}
