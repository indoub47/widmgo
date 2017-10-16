package main

import (
	"encoding/json"
	"strconv"
)

type NullUint8 struct {
	Int   uint8
	Valid bool
}

// UnmarshalJSON unmarshals NullUint8 value from JSON string
func (ni *NullUint8) UnmarshalJSON(b []byte) error {

	if len(b) == 0 {
		*ni = NullUint8{Int: 0, Valid: false}
		return nil
	}

	str, err := strconv.Unquote(string(b))
	if err != nil {
		panic(err)
	}

	ui8, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		panic(err)
	}

	*ni = NullUint8{Int: uint8(ui8), Valid: true}

	return nil
}

func (ni NullUint8) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(ni.Int)
}
