package main

import (
	"strconv"
	"time"
)

type Time struct {
	time.Time
}

// UnmarshalJSON unmarshals Time value from JSON string
func (t *Time) UnmarshalJSON(b []byte) error {
	str, _ := strconv.Unquote(string(b[:]))
	tm, err := time.Parse("2006-01-02", str)
	if err != nil {
		panic(err)
	}
	*t = Time{tm}
	return nil
}
