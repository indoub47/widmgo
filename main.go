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
}
