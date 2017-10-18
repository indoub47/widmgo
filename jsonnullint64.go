package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}

	return json.Marshal(nil)
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *interface{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	// if json null
	if x == nil {
		v.Valid = false
		v.Int64 = 0
		return nil
	}

	// pointer -> to interface{}
	xi := *x
	// if string
	if s, ok := xi.(string); ok {
		// if empty string
		if len(s) == 0 {
			v.Valid = false
			v.Int64 = 0
			return nil
		}

		i, err := strconv.ParseInt(s, 10, 64)

		// if string couldn't be conversed to int64
		if err != nil {
			v.Valid = false
			v.Int64 = 0
			return err
		}

		// if string represents an int64
		v.Valid = true
		v.Int64 = i
		return nil
	}

	switch n := xi.(type) {
	case float64:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	case int:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	case int64:
		v.Valid = true
		v.Int64 = n
		return nil
	case int32:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	case float32:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	case int16:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	case int8:
		v.Valid = true
		v.Int64 = int64(n)
		return nil
	default:
		return errors.New("Bad json field type, unable to convert to JsonNullInt64 type")
	}

}
