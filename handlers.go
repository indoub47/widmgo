package main

import (
	"encoding/json"
	"net/http"

	"github.com/jcelliott/lumber"
)

func save(rw http.ResponseWriter, req *http.Request) {

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
	if err := decoder.Decode(&recs); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		mlog.Error("decoder.Decode() error:", err)
		return
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
		rw.WriteHeader(http.StatusBadRequest)
		js, err := json.Marshal(ves)
		if err != nil {
			mlog.Error("Rec.Validate() error(s), also unsuccessful marshal", err)
		} else {
			mlog.Error("Rec.Validate() error(s):", js)
		}
		return
	}

	// accepted
	//rw.WriteHeader(http.StatusAccepted)

	// mėginti sukišti į db
	_, err = insertRecs(recs)

	// jeigu nepavyksta sukišti į db
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError) // internal error
		mlog.Fatal("Failed insert\n", err)
		if err := json.NewEncoder(rw).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// jeigu pavyko db operacija
	rw.WriteHeader(http.StatusNoContent)
	mlog.Info("Inserted:\n", recs)
}

func receive(rw http.ResponseWriter, req *http.Request) {

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

	mlog.Info("handlers.receive - before fetching")

	// mėginti partempti iš db
	recs, err := fetchRecs()
	if err != nil {
		mlog.Fatal("Failed fetch", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	mlog.Info("handlers.receive - between fetching and writing response")

	// mėginti išsiųsti
	if err := json.NewEncoder(rw).Encode(recs); err != nil {
		mlog.Fatal("Failed encode to response writer", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	mlog.Info("handlers.receive - between writing response and updating db")

	_, err = markAsSent()
	if err != nil {
		mlog.Error("Failed to update sent records in DB", err)
	}

}
