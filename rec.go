package main

import (
	"time"
)

type Rec struct {
	ID          JsonNullInt64
	Linija      string
	Kelias      int64 `json:",string"`
	Km          int64 `json:",string"`
	Pk          int64 `json:",string"`
	M           int64 `json:",string"`
	Siule       JsonNullInt64
	Skodas      string
	Suvirino    JsonNullString
	Operatorius string
	Aparatas    string
	TData       Time
	Kelintas    int64 `json:",string"`
}

func (r *Rec) toSQLArgs() []interface{} {
	var a = make([]interface{}, 14)

	a[0] = r.ID
	a[1] = r.Linija
	a[2] = r.Kelias
	a[3] = r.Km
	a[4] = r.Pk
	a[5] = r.M
	a[6] = r.Siule
	a[7] = r.Skodas
	a[8] = r.Suvirino
	a[9] = r.Operatorius
	a[10] = r.Aparatas
	a[11] = r.TData.Format("2006-01-02")
	a[12] = r.Kelintas
	a[13] = time.Now().Format("2006-01-02T15:04:05")

	return a
}
