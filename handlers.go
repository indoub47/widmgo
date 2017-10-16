package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jcelliott/lumber"
)

func test(rw http.ResponseWriter, req *http.Request) {

	flog, err := lumber.NewFileLogger("log.log", lumber.TRACE, lumber.ROTATE, 5000, 9, 100)
	if err != nil {
		panic(err)
	}
	flog.Prefix("handlers.go-test")
	clog := lumber.NewConsoleLogger(lumber.TRACE)
	mlog := lumber.NewMultiLogger()
	mlog.AddLoggers(flog, clog)

	defer mlog.Close()

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var recs []Rec

	// jeigu nepavyksta dekodinti json:
	// 1. siunčia 400 header
	// 2. panic
	if err := decoder.Decode(&recs); err != nil {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.WriteHeader(http.StatusBadRequest)
		mlog.Error("decoder.Decode() error:", err)
	}

	// validatinti
	var ves = []ValidatedEntity{}
	for _, rec := range recs {
		if ve := rec.Validate(); len(ve.Errors) > 0 {
			ves = append(ves, ve)
		}
	}

	// jeigu buvo blogų, rašyti atsakymą
	if len(ves) > 0 {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.WriteHeader(http.StatusBadRequest)

		js, err := json.Marshal(ves)
		if err != nil {
			mlog.Error("Rec.Validate() error(s), also unsuccessful marshal", err)
		} else {
			mlog.Error("Rec.Validate() error(s):", js)
		}
	}

	// mėginti sukišti į db
	c, err := insertRecs(recs)

	// jeigu nepavyksta sukišti į db
	// 1. siunčia 500 header
	// 2. panic
	if err != nil {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.WriteHeader(http.StatusInternalServerError) // internal error
		if err := json.NewEncoder(rw).Encode(err); err != nil {
			panic(err)
		}
	}

	// jeigu pavyko db operacija
	// 1. siunčia 200 header
	// 2. log
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(http.StatusCreated)

	fmt.Println(recs)
}
