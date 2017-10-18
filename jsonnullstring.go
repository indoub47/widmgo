package main

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}

	return json.Marshal(nil)
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *interface{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	// if json null
	if x == nil {
		v.Valid = false
		v.String = ""
		return nil
	}

	// pointer -> to interface{}
	xi := *x

	if s, ok := xi.(string); ok {
		v.Valid = true
		v.String = s
		return nil
	}

	// if cannot convert to string
	return errors.New("Bad json field type, unable to convert to JsonNullString type")
}
