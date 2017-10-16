package main

type Rec struct {
	Id          uint32 `json:",string,omitempty"`
	Linija      string
	Kelias      uint8  `json:",string"`
	Km          uint16 `json:",string"`
	Pk          uint16 `json:",string"`
	M           uint16 `json:",string"`
	Siule       NullUint8
	Skodas      string
	Suvirino    string
	Operatorius string
	Aparatas    string
	TData       Time
	Kelintas    uint8 `json:",string"`
}
