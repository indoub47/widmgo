package main

import (
	"log"
	"net/http"
)

func main() {
	/* 	router.Handle("/defektai", jwtMiddleware.Handler(http.HandlerFunc(insertHandler))).Methods("POST")
	   	router := mux.NewRouter()
		   log.Fatal(http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router))) */

	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":3000", nil))

	/*
		var r1 = Rec{JsonNullInt64{sql.NullInt64{Int64: 0, Valid: false}}, "17", 1, 142, 10, 23, JsonNullInt64{sql.NullInt64{Int64: 0, Valid: true}}, "06.4", JsonNullString{sql.NullString{String: "gamykla", Valid: true}}, "426", "830", Time{time.Now()}, 1}
		var r2 = Rec{JsonNullInt64{sql.NullInt64{Int64: 0, Valid: false}}, "17", 1, 142, 10, 23, JsonNullInt64{sql.NullInt64{Int64: 9, Valid: true}}, "06.4", JsonNullString{sql.NullString{String: "gamykla", Valid: true}}, "426", "830", Time{time.Now()}, 1}
		var r3 = Rec{JsonNullInt64{sql.NullInt64{Int64: 456, Valid: true}}, "01", 1, 333, 5, 23, JsonNullInt64{sql.NullInt64{Int64: 0, Valid: true}}, "06.4", JsonNullString{sql.NullString{String: "IF4", Valid: true}}, "422", "831", Time{time.Now()}, 2}

		var recs = []Rec{r1, r2, r3}

		var ves = []ValidatedEntity{}
		for _, rec := range recs {
			if ve := rec.Validate(); len(ve.Errors) > 0 {
				ves = append(ves, ve)
			}
		}

		for _, v := range ves {
			fmt.Printf("%+v\n", v)
		}*/

}
